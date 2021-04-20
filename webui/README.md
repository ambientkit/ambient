# Webui

```bash
# Create the NPM file.
npm init

# https://esbuild.github.io/getting-started/#build-scripts
# https://www.codementor.io/@uokesita/build-a-react-js-application-with-esbuild-and-node-1fjklerh4f

# Add the dependencies.
npm install -S esbuild
npm install -S fs-extra
npm install -S chokidar
npm install -S react react-dom

./node_modules/.bin/esbuild src/index.jsx --bundle --outfile=dist/out.js

npx servor ./build/ index.html --reload
```