var express = require('express');
var app     = module.exports = express.createServer();

// Settings
app.set('views', __dirname + '/views');
app.set('view engine', 'jade');
app.set('view cache', !!parseInt(process.env.VIEW_CACHE));
app.set('view options', {
  layout: false
});

// Middleware
app.use(express.bodyParser());
app.use(express.methodOverride());
app.use(app.router);
app.use(express.compiler({src: __dirname + '/public', enable: ['less']}));
app.use(express.static(__dirname + '/public'));
app.use(express.errorHandler({dumpExceptions: true, showStack: true}));

// Helpers
app.helpers(require('./views/helpers'));

// Repositories
var repositories    = require('./repositories');
var talksRepository = new repositories.TalksRepository();

// Controllers
var controllers     = require('./controllers')
var pagesController = new controllers.PagesController({
  talks: talksRepository,
});

// Routes
app.get('/', pagesController.action('homepage'));
app.get('/consulting', pagesController.action('consulting'));
app.get('/speaking', pagesController.action('speaking'));
