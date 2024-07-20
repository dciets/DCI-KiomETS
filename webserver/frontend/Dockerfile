FROM node:16.15.0 as build

WORKDIR /usr/local/app

COPY ./package.json /usr/local/app/

RUN npm install -g @angular/cli
RUN npm install

COPY ./ /usr/local/app/

RUN npm run build

FROM nginx:latest

COPY --from=build /usr/local/app/dist/frontend/ /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf
