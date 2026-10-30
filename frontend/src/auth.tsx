import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
  type ReactNode,
} from "react";

export interface User {
  userID: string;
  username: string;
  globalName: string;
  avatarUrl: string;
  expiresAt: string;
}

type AuthContextValue = {
  user: User | null;
  loading: boolean;
  loggingOut: boolean;
  error: string | null;
  refresh: () => Promise<void>;
  logout: () => Promise<void>;
};

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const [loggingOut, setLoggingOut] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const currentRequest = useRef<AbortController | null>(null);
  const mutatingSession = useRef(false);

  const refresh = useCallback(async () => {
    if (mutatingSession.current) return;
    currentRequest.current?.abort();
    const controller = new AbortController();
    currentRequest.current = controller;
    const timeout = window.setTimeout(() => controller.abort("timeout"), 8000);
    try {
      const response = await fetch("/api/me", {
        credentials: "include",
        signal: controller.signal,
      });
      if (controller !== currentRequest.current) return;
      if (response.status === 401) {
        setUser(null);
        setError(null);
        return;
      }
      if (!response.ok) throw new Error("Session unavailable");
      const data: User = await response.json();
      if (controller !== currentRequest.current || controller.signal.aborted)
        return;
      if (!data.userID || !Number.isFinite(Date.parse(data.expiresAt)))
        throw new Error("Invalid session");
      setUser(data);
      setError(null);
    } catch {
      if (
        controller === currentRequest.current &&
        (!controller.signal.aborted || controller.signal.reason === "timeout")
      ) {
        setError(
          "We couldn’t reach the login service. You can still explore all players.",
        );
      }
    } finally {
      window.clearTimeout(timeout);
      if (controller === currentRequest.current) setLoading(false);
    }
  }, []);

  const logout = useCallback(async () => {
    if (mutatingSession.current) return;
    mutatingSession.current = true;
    setLoggingOut(true);
    // Prevent an older session request from restoring a user after logout.
    currentRequest.current?.abort();
    currentRequest.current = null;
    try {
      const response = await fetch("/auth/logout", {
        method: "POST",
        credentials: "include",
        signal: AbortSignal.timeout(8000),
      });
      if (!response.ok) throw new Error("Logout failed");
      setUser(null);
      setError(null);
    } catch {
      setError("We couldn’t sign you out. Please try again.");
    } finally {
      mutatingSession.current = false;
      setLoggingOut(false);
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    void refresh();
    const onFocus = () => {
      void refresh();
    };
    window.addEventListener("focus", onFocus);
    const interval = window.setInterval(onFocus, 60_000);
    return () => {
      window.removeEventListener("focus", onFocus);
      window.clearInterval(interval);
      currentRequest.current?.abort();
      currentRequest.current = null;
    };
  }, [refresh]);

  useEffect(() => {
    if (!user) return;
    const timer = window.setTimeout(
      () => setUser(null),
      Math.max(0, Date.parse(user.expiresAt) - Date.now()),
    );
    return () => window.clearTimeout(timer);
  }, [user]);

  return (
    <AuthContext.Provider
      value={{ user, loading, loggingOut, error, refresh, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used inside AuthProvider");
  return context;
}
