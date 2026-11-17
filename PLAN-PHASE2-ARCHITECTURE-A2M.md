# PLAN PHASE 2 — ARCHITECTURE A2M TUN

- **Date** : 15 novembre 2026
- **Statut** : proposition — **aucune modification de code avant validation**

---

## 1. Architecture cible (3 couches)

```
A2M TUN UI (Compose, identité A2M)
        │  (collecte d'état + actions)
        ▼
A2M UI State / ViewModels (états réels : DISCONNECTED/CONNECTING/CONNECTED/DISCONNECTING/ERROR)
        │  (réutilise DashboardViewModel + LogViewModel + CommandClient)
        ▼
Moteur existant : BoxService → VPNService (VpnService Android) → libbox (sing-box)
```

- Aucune réécriture du moteur : `VPNService`, `BoxService`, `ServiceBinder`,
  `ServiceNotification`, profils DB, import handler et stats sont réutilisés tels quels.
- L'UI A2M ne fabrique **jamais** un état : elle lit `DashboardViewModel.uiState`
  (alimenté par le vrai service) et `Status`.

## 2. Nouvelle navigation A2M

### Routes A2M (nouvelles)
| Route | Écran | Contenu |
|---|---|---|
| `a2m/home` | Accueil | Header A2M, carte monde, bouton CONNECT, carte serveur, stats, IP |
| `a2m/servers` | Serveurs | Liste des profils réels, import (fichier/lien/presse-papiers), sélection, édition, suppression |
| `a2m/routing` | Routage | Groupes/outbound, proxy par application, DNS, IPv4/IPv6 (fonctions réellement disponibles) |
| `a2m/journal` | Journal | Logs réels (`LogViewModel`) : filtres, recherche, effacement |
| `a2m/settings` | Réglages | Auto-démarrage, reconnexion, notifications, mode sombre, versions app/moteur |
| `a2m/server/{id}` | Détail/édition serveur | Édition du profil si le format le permet |

### Bottom bar (4 onglets, conforme mission)
`Accueil` · `Serveurs` · `Routage` · `Réglages` — `Journal` accessible via le menu hamburger
et/ou un accès secondaire (exigence §13).

### Suppression de l'accès utilisateur aux écrans Sing-box
1. Retirer du `NavHost` les destinations Sing-box devenues inutiles **ou** les rendre
   inaccessibles (aucun point d'entrée UI). Priorité : suppression des routes UI
   Sing-box du graphe, conservation des écrans réutilisés (`LogScreen` peut être
   réhabillé A2M, `ProfilesCard` réutilisé dans la page Serveurs).
2. Remplacer les mappings de la bottom bar et du menu hamburger (actuellement
   `dashboard` / `log` / `settings`) par les routes `a2m/*`.
3. Remplacer `showProfilePickerSheet()` (feuille Sing-box) par la page Serveurs A2M.
4. Supprimer les entrées UI de mise à jour GitHub (sollicitation 1er lancement,
   dialog, badge) sans supprimer le code de mise à jour interne.
5. `app_name` → « A2M TUN » ; libellé d'activité ; icône conservée ou remplacée.
6. Masquer la `ServiceStatusBar` Sing-box sur `a2m/home` (l'accueil A2M affiche déjà
   l'état du tunnel).

## 3. Correspondance exigences → implémentation

| Exigence | Implémentation |
|---|---|
| Bouton CONNECT réel | `dashboardViewModel.toggleService()` + vérification `Status` réel ; libellés Déconnecté/Connexion/Connecté/Déconnexion/Erreur |
| Carte serveur | lit `uiState.selectedProfileName/Id` ; clic → `a2m/servers` |
| Stats | `uiState.downlink/uplink` + totaux (réels) |
| IP sortie | requête réelle (ipify) après connexion, jamais simulée |
| Serveurs | Profils existants + `ProfileImportHandler` (fichier, lien) + `ClipboardImporter` ; validation avant enregistrement |
| Routage | Groupes (`GroupsCard`/`GroupsViewModel`), proxy par app (`PerAppProxyScreen` réutilisé), DNS/IPv6 selon capacités réelles du moteur |
| Journal | `LogViewModel` (niveaux, recherche, `clearLogs`) |
| Réglages | `Settings` existants (auto-start service, notifications, mode sombre) + versions réelles |

## 4. Découpage des phases suivantes

- **Phase 3 — UI** : créer `A2MTUNNavigation` (routes `a2m/*`), `A2MTUNServersScreen`,
  `A2MTUNRoutingScreen`, `A2MTUNSettingsScreen`, `A2MTUNJournalScreen` ; brancher la
  bottom bar et le menu ; retirer les points d'entrée Sing-box ; nom « A2M TUN ».
- **Phase 4 — Moteur** : connexion réelle du bouton, sélection de profil, stats, IP,
  états réels, erreurs ; suppression de la feuille de picker Sing-box.
- **Phase 5 — Validation** : tests CI + vérifications runtime (logcat si dispo).
- **Phase 6 — Compilation** : APK release via CI, vérification de l'artefact.

## 5. Risques et précautions

| Risque | Précaution |
|---|---|
| Casser le moteur | Aucune modification de `bg/*`, `database/*`, profils, import |
| Écrans Sing-box encore atteignables | Retirer les routes du graphe + mappings + deep links |
| Mises à jour GitHub réapparues | Supprimer les déclencheurs UI, garder le code interne inactif |
| Régression Compose | Validation par CI (`:app:compileOtherDebugKotlin` + `assembleOtherRelease`) |
| WebView | Non retenue : Compose natif (déjà en place) |

## 6. Critères d'acceptation

1. Au lancement (y compris à froid), l'écran A2M TUN s'affiche — aucun Dashboard Sing-box.
2. Aucun accès UI aux écrans Sing-box (bottom bar, menu, cartes).
3. Aucune fenêtre de mise à jour GitHub.
4. Les 4 onglets fonctionnent avec des données réelles ; Journal accessible.
5. CONNECT démarre le vrai moteur ; CONNECTÉ uniquement après validation réelle.
6. Import fichier/lien/presse-papiers → profil visible et sélectionnable.
7. Build CI vert + APK installable.
