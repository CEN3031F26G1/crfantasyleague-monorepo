# Catalog sources and image credits

Research snapshot: **October 30, 2026**. Biographies are short, original summaries of the linked public information. This snapshot is not automatically refreshed.

## Tournament and qualification

- [CRL 2026 World Finals — Liquipedia](https://liquipedia.net/clashroyale/Clash_Royale_League/2026/World_Finals#Participants), page revision **160270**, retrieved through the MediaWiki parse API.
- Listed field: 16 players; Shanghai, November 6–8, 2026.
- Five season winners, seven season leaderboard qualifiers, two Last Chance Qualifier players, and two Champions Cup players.
- This is a secondary-source roster, not an official Supercell data integration. Recheck before launch.

## Players

| Player        | Country            | Qualification         | Biography source                                                                                    | Profile image source                                                                                            |
| ------------- | ------------------ | --------------------- | --------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------- |
| Mohamed Light | Egypt              | Season 1 winner       | [Liquipedia](https://liquipedia.net/clashroyale/Mohamed_Light)                                      | [Public profile](https://www.youtube.com/channel/UCjFrTvMA5zvCnb9EAX9MNVw) → `public/players/mohamed-light.jpg` |
| Pedro         | Brazil             | Season 2 winner       | [Liquipedia](https://liquipedia.net/clashroyale/Pedro)                                              | [Public profile](https://www.youtube.com/channel/UCO8wTrnhGkcEhDAhVNB3-ag) → `public/players/pedro.jpg`         |
| SandBox       | South Korea        | Season 3 winner       | [Liquipedia](https://liquipedia.net/clashroyale/SandBox)                                            | [Public profile](https://www.youtube.com/channel/UCUn9lldIkbsDE_DK55PlJMQ) → `public/players/sandbox.jpg`       |
| Clown         | South Korea        | Season 4 winner       | [Liquipedia](https://liquipedia.net/clashroyale/Clown)                                              | Not available — monogram fallback                                                                               |
| Mugi          | Japan              | Season 5 winner       | [Liquipedia](https://liquipedia.net/clashroyale/Mugi)                                               | [Public profile](https://www.youtube.com/channel/UClWUQjL966gilJUUHBZ5-4g) → `public/players/mugi.jpg`          |
| Sub           | Brazil             | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/Sub)                                                | Not available — monogram fallback                                                                               |
| adriel        | Dominican Republic | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/Adriel)                                             | [Public profile](https://x.com/AdrielEduardo11) → `public/players/adriel.jpg`                                   |
| Ryley         | Canada             | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/Ryley)                                              | [Public profile](https://www.youtube.com/channel/UCAw2ZHAIP177ryD3WF9iXmw) → `public/players/ryley.jpg`         |
| Vitor75       | Portugal           | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/KING_VITOR)                                         | Not available — monogram fallback                                                                               |
| Guriko        | Japan              | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/Griko)                                              | [Public profile](https://www.youtube.com/channel/UCqvB8-10UkiTvqnpGK4Fz4A) → `public/players/guriko.jpg`        |
| Morten        | Germany            | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/Morten)                                             | [Public profile](https://www.youtube.com/channel/UCn2AA9OCYnZZ1T82RaGUYPw) → `public/players/morten.jpg`        |
| Viiper        | France             | Season leaderboard    | [Liquipedia](https://liquipedia.net/clashroyale/Viiper)                                             | [Public profile](https://www.youtube.com/channel/UCymILwHUZt5XiC-AN2y-8cA) → `public/players/viiper.jpg`        |
| Woo           | South Korea        | Last chance qualifier | [Liquipedia](https://liquipedia.net/clashroyale/Clash_Royale_League/2026/World_Finals#Participants) | Not available — monogram fallback                                                                               |
| JorZ          | Spain              | Last chance qualifier | [Liquipedia](https://liquipedia.net/clashroyale/JorZ)                                               | Not available — monogram fallback                                                                               |
| Hermit        | China              | Champions Cup         | [Liquipedia](https://liquipedia.net/clashroyale/Clash_Royale_League/2026/World_Finals#Participants) | Not available — monogram fallback                                                                               |
| Midnight      | China              | Champions Cup         | [Liquipedia](https://liquipedia.net/clashroyale/Clash_Royale_League/2026/World_Finals#Participants) | Not available — monogram fallback                                                                               |

## Image provenance and limitations

- Eight profile images were retrieved from the Open Graph image of each player’s YouTube channel. Channel links were taken from their individual Liquipedia profile, rather than guessed from a name. These are the public channel avatars and may be logos/artwork, not photographic headshots.
- adriel’s public X account is linked from his Liquipedia profile. The account avatar URL was obtained from the public FxTwitter API and the image was downloaded from `pbs.twimg.com`.
- Liquipedia had photographic file references for Pedro and SandBox, but the image host returned HTTP 429. Those files are **not** included; their YouTube profile images are used instead.
- No verified usable image was obtained for **Clown, Sub, Vitor75, Woo, JorZ, Hermit, or Midnight**. Their monograms intentionally indicate missing assets, not a verified photo.
- Woo, Hermit, and Midnight did not have individual biography pages in the retrieved roster. Their descriptions only state the supported country and qualification path. No personal history is invented.
- Downloaded profile images are served locally to avoid runtime requests to third-party image hosts. Discord avatars for signed-in users still come from Discord’s CDN.
- Third-party profile images, logos, and artwork remain the property of their respective creators/rights holders. Source attribution is not a license grant. Obtain suitable permissions or replace with approved media before public/commercial deployment.

## Other references

- [RoyaleTools](https://www.royaletools.com/) was inspected for general UI direction: compact navigation, dark surfaces, tool cards, and public access. This implementation uses its own layout, CSS, wordmark, and icons.
- [Discord OAuth2 documentation](https://docs.discord.com/developers/topics/oauth2) describes the standard authorization-code/identify flow used by the existing backend.
