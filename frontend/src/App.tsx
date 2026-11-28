import { useEffect, useState } from "react";
import {
  Link,
  NavLink,
  Navigate,
  Route,
  Routes,
  useLocation,
  useSearchParams,
} from "react-router-dom";
import * as Avatar from "@radix-ui/react-avatar";
import * as DropdownMenu from "@radix-ui/react-dropdown-menu";
import {
  ArrowRight,
  ChevronDown,
  CircleHelp,
  Crown,
  ExternalLink,
  Users,
  LogOut,
  Search,
  ShieldCheck,
  Trophy,
  X,
} from "lucide-react";
import "./arena.css";
import { useAuth } from "./auth";
import {
  DiscordIcon,
  ExternalAnchor,
  FilterSelect,
  Flag,
  initials,
  Logo,
  PlayerDialog,
  PlayerPortrait,
  players,
  rosterSource,
  verifiedDate,
  type Player,
} from "./components";

function AccountMenu() {
  const { user, loading, logout, loggingOut } = useAuth();
  if (loading)
    return (
      <span className="session-loading" role="status">
        Checking session…
      </span>
    );
  if (!user)
    return (
      <Link
        to="/login"
        className="button button-small button-discord"
        aria-label="Sign in with Discord"
      >
        <DiscordIcon size={17} />
        <span>
          Sign in<span className="desktop-only"> with Discord</span>
        </span>
      </Link>
    );
  const name = user.globalName || user.username || "Royale fan";
  return (
    <DropdownMenu.Root>
      <DropdownMenu.Trigger
        className="account-trigger"
        disabled={loggingOut}
        aria-label={`Account menu for ${name}`}
      >
        <Avatar.Root className="account-avatar">
          <Avatar.Image src={user.avatarUrl || undefined} alt="" />
          <Avatar.Fallback>{initials(name)}</Avatar.Fallback>
        </Avatar.Root>
        <span className="account-name">{name}</span>
        <ChevronDown size={14} />
      </DropdownMenu.Trigger>
      <DropdownMenu.Portal>
        <DropdownMenu.Content
          className="account-menu"
          sideOffset={10}
          align="end"
        >
          <DropdownMenu.Label className="account-label">
            Signed in as<strong>{name}</strong>
            {user.username && <span>@{user.username}</span>}
          </DropdownMenu.Label>
          <DropdownMenu.Separator className="menu-separator" />
          <DropdownMenu.Item
            className="menu-item"
            disabled={loggingOut}
            onSelect={() => {
              void logout();
            }}
          >
            <LogOut size={16} />
            {loggingOut ? "Signing out…" : "Sign out"}
          </DropdownMenu.Item>
        </DropdownMenu.Content>
      </DropdownMenu.Portal>
    </DropdownMenu.Root>
  );
}

function Header() {
  const { pathname } = useLocation();
  const arena = pathname === "/" || pathname === "/players" || pathname === "/login";
  return (
    <header className={arena ? "site-header arena-header" : "site-header"}>
      <div className="header-inner">
        <Link to="/" aria-label="Royale Fantasy home">
          <Logo />
        </Link>
        <nav aria-label="Main navigation">
          <NavLink
            to="/"
            end
            className={({ isActive }) =>
              isActive ? "nav-link active" : "nav-link"
            }
          >
            Overview
          </NavLink>
          <NavLink
            to="/players"
            className={({ isActive }) =>
              isActive ? "nav-link active" : "nav-link"
            }
          >
            Players<span className="nav-count">16</span>
          </NavLink>
        </nav>
        <div className="header-account">
          <AccountMenu />
        </div>
      </div>
    </header>
  );
}

function Footer() {
  return (
    <footer className="site-footer">
      <div className="container footer-inner">
        <Link to="/" aria-label="Royale Fantasy home">
          <Logo />
        </Link>
        <small>
          Fantasy only · No real-money wagering · Unaffiliated with Supercell or Discord
        </small>
      </div>
    </footer>
  );
}

function SessionNotice() {
  const { error, refresh, loggingOut, logout } = useAuth();
  const { pathname } = useLocation();
  if (!error || (error === "session" && pathname !== "/login")) return null;
  const failedLogout = error === "logout";
  return (
    <div className="session-notice" role="alert">
      <span>{failedLogout ? "We couldn’t sign you out. Please try again." : "We couldn’t reach the login service."}</span>
      <button
        onClick={() => {
          void (failedLogout ? logout() : refresh());
        }}
        disabled={loggingOut}
      >
        Try again
      </button>
    </div>
  );
}

function Home() {
  return (
    <>
      <section className="arena-hero" aria-labelledby="arena-title">
        <div className="arena-inner">
          <h1 id="arena-title">Your roster. <br /><span>Their moment.</span></h1>
          <p className="arena-intro">
            Sixteen finalists. One crown. Scout the pros on the road to Shanghai
            and find the players you’ll be cheering for.
          </p>
          <div className="arena-actions">
            <Link to="/players" className="arena-cta">
              Explore the players <ArrowRight size={20} aria-hidden="true" />
            </Link>
          </div>
        </div>
        <section className="finals-strip" aria-labelledby="finals-title">
          <div className="finals-identity">
            <span className="finals-crest" aria-hidden="true"><Trophy size={26} /></span>
            <div>
              <p>Clash Royale League · 2026</p>
              <h2 id="finals-title">World Finals</h2>
            </div>
          </div>
          <dl className="finals-facts">
            <div><dt>Finalists</dt><dd>{players.length}</dd></div>
            <div><dt>Dates</dt><dd>6–8 Nov</dd></div>
            <div><dt>Host city</dt><dd>Shanghai</dd></div>
          </dl>
          <ExternalAnchor href={rosterSource} className="finals-details">Event details</ExternalAnchor>
        </section>
      </section>
      <section className="arena-guide" id="how-it-works" aria-labelledby="guide-title">
        <div className="container">
          <div className="arena-divider" aria-hidden="true"><Crown size={24} /></div>
          <h2 id="guide-title">How fantasy works</h2>
          <h2 className="guide-intro">A new way to watch the finals.</h2>
          <div className="how-grid">
            <article>
              <span className="step-icon">
                <Search size={21} aria-hidden="true" />
              </span>
              <span className="step-number">01</span>
              <h3>Know your contenders</h3>
              <p>
                Explore every qualified player, their background, and their road
                to the World Finals.
              </p>
              <Link to="/players" className="inline-link">
                Browse the catalog
                <ArrowRight size={14} aria-hidden="true" />
              </Link>
            </article>
            <article>
              <span className="step-icon">
                <Users size={21} aria-hidden="true" />
              </span>
              <span className="step-number">02</span>
              <h3>Make it your team</h3>
              <p>
                Pick the players you believe in. Put your Clash Royale knowledge
                behind a roster of your own.
              </p>
              <span className="coming-label">Coming soon</span>
            </article>
            <article>
              <span className="step-icon">
                <Crown size={21} aria-hidden="true" />
              </span>
              <span className="step-number">03</span>
              <h3>Follow every crown</h3>
              <p>
                Cheer for your picks on the world stage. Fantasy scoring and
                league competition are on the way.
              </p>
              <span className="coming-label">Coming soon</span>
            </article>
          </div>
        </div>
      </section>
    </>
  );
}

const routeOptions = [
  { value: "all", label: "All qualification routes" },
  { value: "season", label: "Season winners" },
  { value: "leaderboard", label: "Season leaderboard" },
  { value: "last-chance", label: "Last chance qualifier" },
  { value: "champions-cup", label: "Champions Cup" },
];

function PlayerCard({ player }: { player: Player }) {
  return (
    <article className="player-card">
      <div className="player-card-top">
        <PlayerPortrait player={player} />
      </div>
      <div className="player-card-body">
        <div className="player-card-heading">
          <h2>{player.name}</h2>
        </div>
        <p className="player-country"><Flag player={player} /> {player.country}</p>
        <div className="player-route">
          <Trophy size={14} aria-hidden="true" />
          {player.qualification}
        </div>
        <PlayerDialog player={player}>
          <button className="player-details-button" aria-label={`Player profile: ${player.name}`}>
            Player profile
            <ArrowRight size={15} />
          </button>
        </PlayerDialog>
      </div>
    </article>
  );
}

function Catalog() {
  const [params, setParams] = useSearchParams();
  const query = params.get("q") || "";
  const rawRoute = params.get("route") || "all";
  const route = routeOptions.some((option) => option.value === rawRoute)
    ? rawRoute
    : "all";
  const filtered = players.filter(
    (player) =>
      (route === "all" || player.route === route) &&
      `${player.name} ${player.country}`
        .toLowerCase()
        .includes(query.trim().toLowerCase()),
  );
  function updateParam(key: string, value: string) {
    setParams(
      (previous) => {
        const next = new URLSearchParams(previous);
        if (!value || value === "all") next.delete(key);
        else next.set(key, value);
        return next;
      },
      { replace: true },
    );
  }
  return (
    <div className="catalog-page">
      <header className="catalog-banner">
        <div className="container">
          <h1>Meet the <span>finalists.</span></h1>
          <p>Explore the players and discover how they earned their place in the arena.</p>
        </div>
      </header>
      <div className="container catalog-field">
        <div className="catalog-toolbar">
          <div className="search-field">
            <Search size={18} aria-hidden="true" />
            <input
              type="search"
              aria-label="Search players or countries"
              placeholder="Search a player or country…"
              value={query}
              onChange={(event) => updateParam("q", event.target.value)}
            />
            {query && (
              <button
                className="icon-button"
                aria-label="Clear search"
                onClick={() => updateParam("q", "")}
              >
                <X size={15} />
              </button>
            )}
          </div>
          <FilterSelect
            value={route}
            onChange={(value) => updateParam("route", value)}
            options={routeOptions}
          />
        </div>
        <div className="catalog-results-line">
          <span role="status">
            Showing <strong>{filtered.length}</strong> of 16 players
          </span>
        </div>
        {filtered.length ? (
          <div className="player-grid">
            {filtered.map((player) => (
              <PlayerCard key={player.id} player={player} />
            ))}
          </div>
        ) : (
          <div className="empty-state">
            <Search size={28} />
            <h2>No players found</h2>
            <p>Try another player, country, or qualification route.</p>
            <button
              className="button button-outline"
              onClick={() => setParams({})}
            >
              Reset filters
            </button>
          </div>
        )}
        <aside className="catalog-source-note" aria-label="Catalog sources">
          <CircleHelp size={17} />
          <div>
            <p>
              Roster and player biographies sourced from{" "}
              <ExternalAnchor href={rosterSource}>Liquipedia</ExternalAnchor> ·
              Checked {verifiedDate}.
            </p>
          </div>
        </aside>
      </div>
    </div>
  );
}

const loginErrors: Record<string, string> = {
  state:
    "Your login request expired or could not be verified. Please start again.",
  cancelled:
    "Discord sign-in was cancelled. No worries — you can try again or browse the players.",
  unavailable:
    "Discord sign-in is temporarily unavailable. Please try again in a moment.",
};

function Login() {
  const { user, loading } = useAuth();
  const [params] = useSearchParams();
  const error = params.get("error");
  const [leaving, setLeaving] = useState(false);
  useEffect(() => {
    const onPageShow = () => setLeaving(false);
    window.addEventListener("pageshow", onPageShow);
    return () => window.removeEventListener("pageshow", onPageShow);
  }, []);
  if (!loading && user) return <Navigate to="/" replace />;
  return (
    <section className="login-page" aria-labelledby="login-title">
      <div className="login-layout">
        <div className="login-card">
          <img className="login-crest" src="/royale-shield.svg" alt="" width="48" height="48" />
          <h1 id="login-title">Sign in to<br /><span>Royale Fantasy.</span></h1>
          <SessionNotice />
          {error && (
            <div className="login-error" role="alert">
              {loginErrors[error] ||
                "Something went wrong while signing in. Please try again."}
            </div>
          )}
          <a
            className={`button discord-login-button ${leaving ? "is-loading" : ""}`}
            href="/auth/login"
            aria-disabled={leaving}
            onClick={(event) => {
              if (leaving) event.preventDefault();
              else setLeaving(true);
            }}
          >
            <DiscordIcon size={20} />
            {leaving ? "Opening Discord…" : "Continue with Discord"}
            <ExternalLink size={15} />
          </a>
          <section className="login-security" aria-labelledby="login-security-title">
            <h2 id="login-security-title">
              <ShieldCheck size={17} aria-hidden="true" />
              Your privacy
            </h2>
            <ul>
              <li>Your password stays with Discord.</li>
              <li>Only your name, avatar and ID are shared.</li>
              <li>Your messages and email stay private.</li>
            </ul>
          </section>
        </div>
      </div>
    </section>
  );
}

function NotFound() {
  return (
    <div className="container not-found">
      <span className="eyebrow gold">404 · WRONG ARENA</span>
      <h1>This page isn’t in the deck.</h1>
      <p>Head back to the overview or scout the finalists.</p>
      <Link to="/" className="button button-gold">
        Back to overview
        <ArrowRight size={16} />
      </Link>
    </div>
  );
}

export default function App() {
  const { pathname } = useLocation();
  useEffect(() => {
    document.title = `${pathname === "/players" ? "Player catalog" : pathname === "/login" ? "Sign in" : "Your roster. Their moment."} — Royale Fantasy`;
    if (pathname === "/" && window.location.hash === "#how-it-works") {
      document.getElementById("how-it-works")?.scrollIntoView();
    } else {
      window.scrollTo(0, 0);
    }
  }, [pathname]);
  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <Header />
      {pathname !== "/login" && <SessionNotice />}
      <main id="main">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/players" element={<Catalog />} />
          <Route path="/login" element={<Login />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      <Footer />
    </>
  );
}
