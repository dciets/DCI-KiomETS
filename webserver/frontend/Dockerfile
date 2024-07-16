FROM node:latest as build

WORKDIR /usr/local/app

COPY ./package.json /usr/local/app/

RUN npm install

COPY ./ /usr/local/app/


RUN npm run build

FROM nginx:latest

COPY --from=build /usr/local/app/dist/frontend/ /usr/share/nginx/html
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
