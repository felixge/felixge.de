var _     = require('underscore');
var async = require('async');

module.exports = PagesController;
function PagesController(properties) {
  this._talks = properties.talks;
}

PagesController.prototype.action = function(name) {
  var action = this[name].bind(this);
  return function(req, res) {
    res.local('breadcrumbs', [{
      title : 'homepage',
      url   : '/'
    }]);
    action(req, res);
  };
};

PagesController.prototype.homepage = function(req, res) {
  async.parallel({
    recentTalks: this._talks.findRecent,
    recentPosts: function(cb) {
      cb(null, []);
    },
  }, function(err, results) {
    res.render('pages/homepage', _.extend(results, {
      title: 'Felix Geisendörfer',
    }));
  });
};

PagesController.prototype.consulting = function(req, res,f) {
  res.render('pages/consulting', {title: 'Consulting'});
};

PagesController.prototype.speaking = function(req, res,f) {
  async.parallel({
    talks: this._talks.findTalks,
  }, function(err, results) {
    res.render('pages/speaking', _.extend(results, {
      title: 'Speaking',
    }));
  });
};
