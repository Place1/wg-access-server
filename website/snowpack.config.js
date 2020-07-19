module.exports = {
  extends: '@snowpack/app-scripts-react',
  scripts: {
    "run:codegen": "npm run codegen"
  },
  plugins: [],
  buildOptions: {
    minify: false,
  },
  devOptions: {
    port: 3000,
    open: 'none'
  },
  proxy: {
    '/api': 'http://localhost:8000/api'
  },
  installOptions: {
    namedExports: ['google-protobuf'],
    rollup: {
      plugins: [require('rollup-plugin-node-polyfills')()],
    }
  }
}
