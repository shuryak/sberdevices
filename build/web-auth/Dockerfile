FROM node:24.4-alpine3.22 AS builder

WORKDIR /web-auth

COPY ./web/auth/package*.json ./
RUN npm install

COPY ./web/auth .

ARG VITE_API_BASE_URL
RUN npm install -g typescript vite
RUN npm run build

FROM busybox:1.37.0-uclibc

COPY --from=builder /web-auth/dist /www

EXPOSE 3000

CMD ["httpd", "-f", "-v", "-p", "3000", "-h", "/www"]
