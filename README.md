# DCI-KiomETS
 
## How to run

`docker compose up` sur la racine du projet si vous utilisez docker.

## Env variables

### Frontend

`/webserver/frontend/src/environments`
Les variables d'environement sont stocké dans un fichier typescript.

- `production` : Laisser tel quel
- `serverAddr` : Adresse du backend

### Backend

`/docker-compose.yml` si vous voulez modifier les variables d'environnement à partir de docker

- `APP_ENV` : Laisser tel quel si vous utilez les variables d'environement par docker, sinon enlever

`/webserver/backend/.env` si vous n'utiliser pas docker / pour faire des test. Les variables du .env sont aussi présentes dans le docker compose

- `GAME_PORT` : Port tcp pour écouter les communications de type « jeu »
- `ADMIN_PORT` : Port tcp pour écouter les communications de type « admin »
- `GAME_SERVER_HOST` : Adresse du serveur de jeu
- `ORIGIN_ALLOWED` : Liste des origines valides pour les requêtes HTTP (venant du frontent)

### Game server

`/docker-compose.yml` si vous voulez modifier les variables d'environnement à partir de docker

- `APP_ENV` : Laisser tel quel si vous utilez les variables d'environement par docker, sinon enlever

`/gameserver/server/.env` si vous n'utiliser pas docker / pour faire des test. Les variables du .env sont aussi présentes dans le docker compose

- `GAME_PORT` : Port tcp pour écouter les communications de type « jeu »
- `ADMIN_PORT` : Port tcp pour écouter les communications de type « admin »
- `SUPER_ADMIN_PORT` : Port tcp pour écouter les communications de type « super admin »
