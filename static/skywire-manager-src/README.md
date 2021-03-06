# Skywire Manager

Frontend application that allows to manage a group of Skywire nodes through a Hypervisor instance.

Note: The software is still under heavy development and the current version is intended for public testing purposes only.

## Prerequisites

The Skywire Manager requires Node 10.9.0 or higher, together with NPM 6.0.0 or higher.

## Initial configuration

Dependencies needed for this app are managed with NPM and are not included in the repository, so you must run the `npm install`
command on this folder before being able to tun the app.

Also, if you are going to use an HTTPS connection, you may want to create a SSL certificate for the dev server. This step is
optional, as the dev server can create its own certificate, but there are 2 reasons for manually creating one: to avoid a bug
in the dev server that sometimes makes it to reload the app when it shouild not
(https://github.com/angular/angular-cli/issues/5826) and to have more freedom managing the trusted certificates. For creating
a custom certificate, follow the steps in the [ssl folder](./ssl/README.md)

## Hypervisor configuration

The Hypervisor instance must be running in `http://127.0.0.1:8000`. If it is running in another URL, you can change it in
[proxy.config.json](proxy.config.json) before running the app.

## Running the app

If the hypervisor instance is running with TLS active (check the hypervisor configuration file) Run `npm run start` to start a
dev server. If you followed the steps indicated in the [ssl folder](./ssl/README.md), the server will use your custom SSL
certificate. If not, the server will use an automatically created one. Alternatively, If the hypervisor instance is running
without TLS, you can start the dev server by running `npm run start-no-ssl`.

After starting the server with `npm run start`, you can access the app by navigating to `https://localhost:4200` with a web
browser (note that you could get a security warning if the SSL certificate is not in the trusted certificates list). Yo can use
`https://localhost:4200` if you started the dev server with `npm run start-no-ssl`. The app will be automatically reloaded if you
change any of the source files.

## Build

Run `npm run build` to build the project. The build artifacts will be stored in the `dist/` directory.

## Translations

You can find information about how to work with translation files in the [Translations README](/static/skywire-manager-src/src/assets/i18n/README.md).
