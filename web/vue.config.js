// vue.config.js
module.exports = {

  // publicPath: process.env.NODE_ENV === 'production' ? '/' : '/dev',
  publicPath: '/',

  // proxy all webpack dev-server requests starting with /api
  // to our backend (localhost:8088) using http-proxy-middleware
  // see https://cli.vuejs.org/config/#devserver-proxy
  devServer: {
    port: 4000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080', // this configuration needs to correspond to the backends' server port
        ws: true,
        changeOrigin: true
      }
    },
  },
  
  // Change build paths to make output 
  // see https://cli.vuejs.org/config/
  outputDir: '../public',
  assetsDir: 'static'
};
