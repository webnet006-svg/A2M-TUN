# PROGRESS — DAMS TUNNEL

## Fait
- Module d'import presse-papiers (`Importer.kt`, `ClipboardImporter.kt`, boutons dans `ProfilesCard.kt`) + 21 tests unitaires.
- Build CI GitHub Actions (branche `dams/ci`) : libbox arm64 + APK debug, tests OK (run 43546409135).
- Build release arm64 unique (variante `other`), signature avec clé de test via secrets GitHub.

## En cours
- Étape 0 : validation réelle sur téléphone (checklist coller / connecter / IP / reconnexion / persistance).

## À faire
- Remplir la checklist étape 0, puis créer les tags `baseline-ok` et `v0-fonctionnel`.
- Phase 1 : thème Dams Cyan (après validation de l'étape 0).

## Problèmes connus
- Le cœur est sur `testing` (1.15 alpha) ; la recommandation 1.13.x n'est pas appliquée.
- `applicationId` et branding encore ceux du fork ; changement prévu avec l'UI.
- La clé de signature de test est dans les secrets GitHub, pas dans le dépôt.
