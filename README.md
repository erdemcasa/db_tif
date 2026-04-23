# db_tif

db_tif est une application web minimaliste qui permet d'enregistrer et de lister des entrées de fichiers TIF. Elle est construite avec un backend en Go, une base de données SQLite3, et un frontend natif en HTML/CSS et Vanilla JS (aucun framework lourd).

## Fonctionnalités

- Interface utilisateur simple et épurée.
- Formulaire d'ajout avec validation (le label doit se terminer par "tif").
- Liste des entrées triées par date de création.
- Traçabilité silencieuse : l'adresse IP et le User-Agent de l'utilisateur sont enregistrés automatiquement dans la base de données lors de la création d'une entrée, sans être exposés sur l'interface publique.

## Prérequis

Pour faire tourner ce projet sur une autre machine, vous aurez besoin de :
- [Go](https://go.dev/doc/install) (version 1.16 ou supérieure recommandée)
- Un compilateur C (comme GCC) car le driver `go-sqlite3` utilise CGO pour fonctionner avec SQLite.

## Installation et démarrage

1. Clonez ce dépôt ou copiez les fichiers sur votre machine.
2. Ouvrez un terminal et placez-vous dans le dossier du projet :
   ```bash
   cd chemin/vers/db_tif
   ```

3. Téléchargez les dépendances du projet (le driver SQLite) :
   ```bash
   go mod download
   ```

4. Lancez le serveur :
   ```bash
   go run main.go
   ```
   *Note : Au premier lancement, le fichier de base de données `db_tif.sqlite` sera créé automatiquement à la racine du projet.*

5. Ouvrez votre navigateur et accédez à :
   ```text
   http://localhost:8080
   ```

## Compilation pour la production

Si vous souhaitez héberger l'application ou la faire tourner en tâche de fond sans utiliser `go run`, vous pouvez compiler un exécutable :

```bash
# Compilation
go build -o db_tif_server main.go

# Lancement de l'exécutable
./db_tif_server
```

Assurez-vous que le dossier `static` reste au même niveau que l'exécutable pour que le serveur puisse trouver les fichiers HTML, CSS et JS.
