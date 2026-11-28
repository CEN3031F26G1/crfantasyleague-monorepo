import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { describe, expect, it, vi } from "vitest";
import App from "../App";
import { AuthProvider } from "../auth";
import { players } from "../components";

function mount(path = "/", signedIn = false) {
  const fetchMock = vi.fn(async (input: RequestInfo | URL) => {
    if (String(input) === "/auth/logout")
      return new Response(null, { status: 204 });
    return signedIn
      ? Response.json({
          userID: "123",
          username: "arenafan",
          globalName: "Arena Fan",
          avatarUrl: "",
          expiresAt: new Date(Date.now() + 3600_000).toISOString(),
        })
      : new Response(null, { status: 401 });
  });
  vi.stubGlobal("fetch", fetchMock);
  render(
    <MemoryRouter initialEntries={[path]}>
      <AuthProvider>
        <App />
      </AuthProvider>
    </MemoryRouter>,
  );
  return fetchMock;
}

describe("public frontend", () => {
  it("shows the landing page and a public player catalog link", async () => {
    mount();
    expect(
      screen.getByRole("heading", { name: /Your roster. Their moment./ }),
    ).toBeInTheDocument();
    expect(
      screen.getByRole("link", { name: /Explore the players/ }),
    ).toHaveAttribute("href", "/players");
    expect(
      await screen.findByRole("link", { name: /Sign in with Discord/ }),
    ).toHaveAttribute("href", "/login");
    const footer = within(screen.getByRole("contentinfo"));
    expect(footer.getAllByRole("link")).toHaveLength(1);
    expect(footer.getByRole("link", { name: "Royale Fantasy home" })).toHaveAttribute("href", "/");
  });

  it("shows all 16 sourced qualifiers without signing in", () => {
    mount("/players");
    expect(
      screen.getAllByRole("button", { name: /Player profile/ }),
    ).toHaveLength(16);
    expect(new Set(players.map((p) => p.id)).size).toBe(16);
    for (const player of players) {
      expect(
        screen.getByRole("heading", { name: player.name }),
      ).toBeInTheDocument();
      expect(player.source.startsWith("https://liquipedia.net/")).toBe(true);
    }
  });

  it("filters by player and country, supports an empty state, and clears filters", async () => {
    mount("/players");
    const user = userEvent.setup();
    const input = screen.getByRole("searchbox");
    await user.type(input, "Japan");
    expect(
      screen.getAllByRole("button", { name: /Player profile/ }),
    ).toHaveLength(2);
    await user.clear(input);
    await user.type(input, "no-such-player");
    expect(
      screen.getByRole("heading", { name: "No players found" }),
    ).toBeInTheDocument();
    await user.click(screen.getByRole("button", { name: "Reset filters" }));
    expect(
      screen.getAllByRole("button", { name: /Player profile/ }),
    ).toHaveLength(16);
  });

  it("supports shareable qualification filters", () => {
    mount("/players?route=champions-cup");
    expect(
      screen.getAllByRole("button", { name: /Player profile/ }),
    ).toHaveLength(2);
    expect(screen.getByRole("heading", { name: "Hermit" })).toBeInTheDocument();
    expect(
      screen.getByRole("heading", { name: "Midnight" }),
    ).toBeInTheDocument();
  });

  it("opens a player profile with initials and sources when no photo is available", async () => {
    mount("/players?q=Hermit");
    const user = userEvent.setup();
    await user.click(screen.getByRole("button", { name: /Player profile/ }));
    const dialog = screen.getByRole("dialog", { name: "Hermit" });
    expect(within(dialog).getByText("H")).toBeInTheDocument();
    expect(within(dialog).queryByText(/A verified profile image is not available/)).toBeNull();
    expect(
      within(dialog).getByRole("link", { name: /Player source/ }),
    ).toHaveAttribute("href", expect.stringContaining("liquipedia.net"));
    await user.keyboard("{Escape}");
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("uses the actual backend OAuth entry point and displays cancellation errors", () => {
    mount("/login?error=cancelled");
    expect(
      screen.getByRole("link", { name: /Continue with Discord/ }),
    ).toHaveAttribute("href", "/auth/login");
    expect(screen.getByRole("alert")).toHaveTextContent(
      "Discord sign-in was cancelled",
    );
    const privacy = within(screen.getByRole("region", { name: "Your privacy" }));
    expect(privacy.getAllByRole("listitem")).toHaveLength(3);
    expect(privacy.getByText("Your password stays with Discord.")).toBeInTheDocument();
  });

  it("renders a not found page", () => {
    mount("/missing");
    expect(
      screen.getByRole("heading", { name: /This page isn’t in the deck/ }),
    ).toBeInTheDocument();
  });
});

describe("authentication UI", () => {
  it("replaces sign in with a profile menu and logs out without gating the catalog", async () => {
    const fetchMock = mount("/players", true);
    const user = userEvent.setup();
    await user.click(
      await screen.findByRole("button", { name: /Account menu for Arena Fan/ }),
    );
    expect(screen.getByText("@arenafan")).toBeInTheDocument();
    await user.click(screen.getByRole("menuitem", { name: /Sign out/ }));
    await waitFor(() =>
      expect(
        screen.getByRole("link", { name: /Sign in with Discord/ }),
      ).toBeInTheDocument(),
    );
    expect(fetchMock).toHaveBeenCalledWith(
      "/auth/logout",
      expect.objectContaining({ method: "POST", credentials: "include" }),
    );
    expect(
      screen.getAllByRole("button", { name: /Player profile/ }),
    ).toHaveLength(16);
  });

  it("redirects an already signed-in visitor away from login", async () => {
    mount("/login", true);
    expect(
      await screen.findByRole("heading", {
        name: /Your roster. Their moment./,
      }),
    ).toBeInTheDocument();
  });

  it("only shows a failed session check when the visitor goes to sign in", async () => {
    const user = userEvent.setup();
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("offline")));
    render(
      <MemoryRouter initialEntries={["/"]}>
        <AuthProvider>
          <App />
        </AuthProvider>
      </MemoryRouter>,
    );
    await screen.findByRole("link", { name: /Sign in with Discord/ });
    expect(screen.queryByRole("alert")).toBeNull();
    await user.click(screen.getByRole("link", { name: /Explore the players/ }));
    expect(
      screen.getAllByRole("button", { name: /Player profile/ }),
    ).toHaveLength(16);
    expect(screen.queryByRole("alert")).toBeNull();
    await user.click(screen.getByRole("link", { name: /Sign in with Discord/ }));
    expect(await screen.findByRole("alert")).toHaveTextContent(
      "We couldn’t reach the login service",
    );
    expect(
      screen.getByRole("button", { name: "Try again" }),
    ).toBeInTheDocument();
  });
});
