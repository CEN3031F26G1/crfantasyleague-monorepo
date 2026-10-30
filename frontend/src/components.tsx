import * as Avatar from "@radix-ui/react-avatar";
import * as Dialog from "@radix-ui/react-dialog";
import * as Select from "@radix-ui/react-select";
import { Check, ChevronDown, Crown, ExternalLink, X } from "lucide-react";
import type { ReactNode } from "react";
import roster from "./data/players.json";

export type Player = (typeof roster)[number];
export const players = roster;
export const rosterSource =
  "https://liquipedia.net/clashroyale/Clash_Royale_League/2026/World_Finals";
export const verifiedDate = "October 30, 2026";

export function Logo({ compact = false }: { compact?: boolean }) {
  return (
    <span className="brand">
      <span className="brand-mark">
        <Crown size={23} strokeWidth={1.8} />
      </span>
      {!compact && (
        <span>
          royale<span className="brand-light">fantasy</span>
          <span className="brand-beta">BETA</span>
        </span>
      )}
    </span>
  );
}

export function DiscordIcon({ size = 20 }: { size?: number }) {
  return (
    <svg
      width={size}
      height={size}
      viewBox="0 0 24 24"
      fill="currentColor"
      aria-hidden="true"
    >
      <path d="M20.32 4.37a19.8 19.8 0 0 0-4.89-1.52l-.6 1.24a18.25 18.25 0 0 0-5.66 0l-.6-1.24a19.7 19.7 0 0 0-4.89 1.52C.59 8.94-.25 13.4.17 17.8a19.9 19.9 0 0 0 6 3.03l1.22-1.98a13 13 0 0 1-1.92-.93l.47-.37c3.7 1.71 7.72 1.71 11.38 0l.47.37a13 13 0 0 1-1.92.93l1.22 1.98a19.9 19.9 0 0 0 6-3.03c.5-5.1-.86-9.52-3.77-13.43ZM8.02 15.13c-1.1 0-2-1-2-2.23 0-1.24.88-2.24 2-2.24 1.11 0 2.01 1.01 2 2.24 0 1.23-.89 2.23-2 2.23Zm7.96 0c-1.1 0-2-1-2-2.23 0-1.24.88-2.24 2-2.24 1.11 0 2.01 1.01 2 2.24 0 1.23-.89 2.23-2 2.23Z" />
    </svg>
  );
}

export function initials(name: string) {
  return name
    .split(/\s+/)
    .slice(0, 2)
    .map((part) => part[0])
    .join("")
    .toUpperCase();
}

export function PlayerPortrait({
  player,
  className = "",
}: {
  player: Player;
  className?: string;
}) {
  return (
    <Avatar.Root className={`player-portrait ${className}`}>
      {player.image && (
        <Avatar.Image
          src={player.image}
          alt={`${player.name} profile image`}
          className="portrait-image"
        />
      )}
      <Avatar.Fallback
        className="portrait-fallback"
        delayMs={player.image ? 200 : 0}
      >
        <span aria-hidden="true" className="portrait-watermark">
          <Crown />
        </span>
        <span>{initials(player.name)}</span>
      </Avatar.Fallback>
    </Avatar.Root>
  );
}

export function Flag({ player }: { player: Player }) {
  const flag = [...player.countryCode]
    .map((letter) => String.fromCodePoint(127397 + letter.charCodeAt(0)))
    .join("");
  return (
    <span role="img" aria-label={player.country} className="flag">
      {flag}
    </span>
  );
}

export function ExternalAnchor({
  href,
  children,
  className = "",
}: {
  href: string;
  children: ReactNode;
  className?: string;
}) {
  return (
    <a
      href={href}
      target="_blank"
      rel="noreferrer"
      className={`external-link ${className}`}
    >
      {children}
      <ExternalLink size={13} aria-hidden="true" />
      <span className="sr-only"> (opens in a new tab)</span>
    </a>
  );
}

export function PlayerDialog({
  player,
  children,
}: {
  player: Player;
  children: ReactNode;
}) {
  return (
    <Dialog.Root>
      <Dialog.Trigger asChild>{children}</Dialog.Trigger>
      <Dialog.Portal>
        <Dialog.Overlay className="dialog-overlay" />
        <Dialog.Content className="dialog-content">
          <Dialog.Close
            className="icon-button dialog-close"
            aria-label="Close player profile"
          >
            <X size={20} />
          </Dialog.Close>
          <PlayerPortrait player={player} className="dialog-portrait" />
          <span className="eyebrow">CRL 2026 · WORLD FINALS</span>
          <Dialog.Title>{player.name}</Dialog.Title>
          <p className="player-country">
            <Flag player={player} /> {player.country}
          </p>
          <Dialog.Description className="dialog-bio">
            {player.bio}
          </Dialog.Description>
          <div className="qualification-detail">
            <span className="muted">Road to Shanghai</span>
            <strong>
              <Check size={16} />
              {player.qualification}
            </strong>
          </div>
          {!player.image && (
            <p className="photo-note">
              A verified profile image is not available yet. A monogram is shown
              instead.
            </p>
          )}
          <div className="dialog-sources">
            <ExternalAnchor href={player.source}>Player source</ExternalAnchor>
            {player.imageSource && (
              <ExternalAnchor href={player.imageSource}>
                Photo source
              </ExternalAnchor>
            )}
          </div>
          <p className="photo-note">
            Roster checked {verifiedDate}. Player information is a sourced
            snapshot, not live data.
          </p>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}

export function FilterSelect({
  value,
  onChange,
  options,
}: {
  value: string;
  onChange: (value: string) => void;
  options: { value: string; label: string }[];
}) {
  return (
    <Select.Root value={value} onValueChange={onChange}>
      <Select.Trigger
        className="filter-select"
        aria-label="Filter by qualification route"
      >
        <Select.Value />
        <Select.Icon>
          <ChevronDown size={16} />
        </Select.Icon>
      </Select.Trigger>
      <Select.Portal>
        <Select.Content
          className="select-content"
          position="popper"
          sideOffset={6}
        >
          <Select.Viewport>
            {options.map((option) => (
              <Select.Item
                key={option.value}
                value={option.value}
                className="select-item"
              >
                <Select.ItemText>{option.label}</Select.ItemText>
                <Select.ItemIndicator>
                  <Check size={14} />
                </Select.ItemIndicator>
              </Select.Item>
            ))}
          </Select.Viewport>
        </Select.Content>
      </Select.Portal>
    </Select.Root>
  );
}
