# PROGRESS — DAMS TUNNEL

## Fait
- Module d'import presse-papiers (`Importer.kt`, `ClipboardImporter.kt`, boutons dans `ProfilesCard.kt`) + 21 tests unitaires.
- Build CI GitHub Actions (branche `dams/ci`) : libbox arm64 + APK debug, tests OK (run 43546409135).
- Build release arm64 unique (variante `other`), signature avec clé de test via secrets GitHub : run 43551265782, APK 35 Mo, artefact 34,2 Mo, tests 21/21, pré-release `dams-v1.15.0-alpha.13-arm64`.
- APK disponible dans `adjuntos/DAMS-TUNNEL-1.15.0-alpha.13-arm64.apk` (ignoré par Git) et envoyé dans le chat.

## En cours
- Étape 0 : validation réelle sur téléphone (checklist coller / connecter / IP / reconnexion / persistance).

## À faire
- Remplir la checklist étape 0, puis créer les tags `baseline-ok` et `v0-fonctionnel`.
- Phase 1 : thème Dams Cyan et branding (après validation de l'étape 0).

## Problèmes connus
- Le cœur est sur `testing` (1.15 alpha) ; la recommandation 1.13.x n'est pas appliquée.
- `applicationId` et label encore ceux du fork (`io.nekohasekai.sfa`, « sing-box ») ; changement prévu avec l'UI.
- La clé de signature de test est dans les secrets GitHub, pas dans le dépôt.
