FROM node:19.0.0-alpine3.16

WORKDIR /workdir/website
COPY . .
RUN npm ci

ENTRYPOINT ["/workdir/website/node_modules/.bin/docusaurus"]
CMD ["-h"]
