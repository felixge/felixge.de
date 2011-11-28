var express = require('express');
var app     = module.exports = express.createServer();

// Settings
app.set('views', __dirname + '/views');
app.set('view engine', 'jade');

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
var controllers        = require('./controllers')
var homepageController = new controllers.HomepageController({
  talks: talksRepository,
});

// Routes
app.get('/', homepageController.index.bind(homepageController));
