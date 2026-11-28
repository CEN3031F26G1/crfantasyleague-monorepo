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
  Swords,
  ArrowLeft,
  ArrowRight,
  Check,
  ChevronDown,
  CircleHelp,
  Crown,
  ExternalLink,
  Layers3,
  LogOut,
  Search,
  ShieldCheck,
  Trophy,
  Users,
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
  const home = useLocation().pathname === "/";
  return (
    <header className={home ? "site-header arena-header" : "site-header"}>
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
          {home && <a href="#how-it-works" className="nav-link">How it works</a>}
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
      <div className="container footer-top">
        <div>
          <Link to="/" aria-label="Royale Fantasy home">
            <Logo />
          </Link>
          <p>A little more at stake. A lot more to cheer for.</p>
        </div>
        <div className="footer-links">
          <Link to="/players">Player catalog</Link>
          <ExternalAnchor href={rosterSource}>Tournament source</ExternalAnchor>
        </div>
      </div>
      <div className="container footer-bottom">
        <span>
          An independent fan project. Not affiliated with Supercell or Discord.
        </span>
        <span>Fantasy only · No real-money wagering</span>
      </div>
    </footer>
  );
}

function SessionNotice() {
  const { error, refresh, loggingOut, logout } = useAuth();
  if (!error) return null;
  const failedLogout = error.includes("sign you out");
  return (
    <div className="session-notice" role="alert">
      <span>{error}</span>
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

const scoutingSteps = [
  { icon: Swords, title: "Meet the pros", detail: "Explore all 16 World Finals contenders." },
  { icon: Crown, title: "Scout your picks", detail: "Discover each player’s path to the finals." },
  { icon: Layers3, title: "Draft your team", detail: "Roster building is coming soon." },
  { icon: Trophy, title: "Compete in leagues", detail: "Fantasy scoring and standings are on the way." },
];

function Home() {
  return (
    <>
      <section className="arena-hero" aria-labelledby="arena-title">
        <div className="arena-inner">
          <p className="arena-eyebrow">Clash Royale League · World Finals 2026</p>
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
      </section>
      <section className="arena-guide" id="how-it-works" aria-labelledby="guide-title">
        <div className="container">
          <div className="arena-divider" aria-hidden="true"><Crown size={24} /></div>
          <h2 id="guide-title">How fantasy works</h2>
          <p className="guide-intro">Get to know the field before you draft your first team.</p>
          <ul className="scouting-steps">
            {scoutingSteps.map(({ icon: Icon, title, detail }) => (
              <li key={title}>
                <span className="arena-emblem"><Icon size={28} aria-hidden="true" /></span>
                <h3>{title}</h3>
                <p>{detail}</p>
              </li>
            ))}
          </ul>
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

function PlayerCard({ player, index }: { player: Player; index: number }) {
  return (
    <article className={`player-card player-color-${index % 5}`}>
      <div className="player-card-top">
        <PlayerPortrait player={player} />
        <span className="qualified-badge">
          <Check size={11} />
          QUALIFIED
        </span>
      </div>
      <div className="player-card-body">
        <div className="player-card-heading">
          <h2>{player.name}</h2>
          <Flag player={player} />
        </div>
        <p className="player-country">{player.country}</p>
        <p className="player-bio">{player.bio}</p>
        <div className="player-route">
          <span className="route-dot" />
          {player.qualification}
        </div>
        <PlayerDialog player={player}>
          <button className="player-details-button">
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
    <div className="container catalog-page">
      <div className="page-breadcrumb">
        <Link to="/">Overview</Link>
        <span>/</span>Player catalog
      </div>
      <div className="catalog-heading">
        <div>
          <span className="eyebrow gold">SHANGHAI · 6–8 NOVEMBER 2026</span>
          <h1>
            Meet the <span>finalists.</span>
          </h1>
          <p>
            The best in the arena. The start of your dream team.
            <br />
            Explore all 16 players qualified for the CRL 2026 World Finals.
          </p>
        </div>
        <div className="catalog-total">
          <Users size={21} />
          <strong>
            16<span>qualified players</span>
          </strong>
        </div>
      </div>
      <div className="catalog-info">
        <ShieldCheck size={17} />
        <span>
          Scouting is open to everyone. Team selection and fantasy scoring are
          coming soon.
        </span>
      </div>
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
        <span>
          <span className="live-dot" />
          2026 WORLD FINALS
        </span>
      </div>
      {filtered.length ? (
        <div className="player-grid">
          {filtered.map((player) => (
            <PlayerCard
              key={player.id}
              player={player}
              index={players.indexOf(player)}
            />
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
          <p>
            This is a manually curated snapshot, not a live feed. Verified
            profile images are shown where available; otherwise, players have a
            monogram. Image credits and individual sources are in each profile.
          </p>
        </div>
      </aside>
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
    <section className="login-page container">
      <Link to="/" className="back-link">
        <ArrowLeft size={15} />
        Back to overview
      </Link>
      <div className="login-layout">
        <div className="login-intro">
          <span className="eyebrow gold">
            WELCOME TO YOUR SIDE OF THE ARENA
          </span>
          <h1>
            Good picks.
            <br />
            <span>Great company.</span>
          </h1>
          <p>
            A fantasy league for people who know that the best matches are the
            ones you’re invested in.
          </p>
          <div className="login-intro-detail">
            <span className="step-icon">
              <Crown size={25} />
            </span>
            <div>
              <strong>One world stage. Your perspective.</strong>
              <span>CRL 2026 World Finals · Shanghai</span>
            </div>
          </div>
        </div>
        <div className="login-card">
          <span className="login-icon">
            <DiscordIcon size={31} />
          </span>
          <h2>Take your seat.</h2>
          <p>
            Sign in with Discord to get your profile ready for Royale Fantasy.
          </p>
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
          <div className="login-security">
            <ShieldCheck size={15} />
            <span>Secure sign-in. No new password.</span>
          </div>
          <div className="login-permissions">
            <p>Just the essentials.</p>
            <span>
              We only request your public Discord profile: your name, avatar,
              and user ID. We never access your messages, servers, or email.
            </span>
          </div>
          <p className="login-footnote">
            Not ready to sign in?{" "}
            <Link to="/players">
              Explore the player catalog
              <ArrowRight size={12} />
            </Link>
          </p>
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
    window.scrollTo(0, 0);
  }, [pathname]);
  return (
    <>
      <a className="skip-link" href="#main">
        Skip to content
      </a>
      <Header />
      {pathname !== "/" && <SessionNotice />}
      <main id="main">
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/players" element={<Catalog />} />
          <Route path="/login" element={<Login />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </main>
      {pathname === "/" && <SessionNotice />}
      <Footer />
    </>
  );
}
