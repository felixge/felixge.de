var _     = require('underscore');
var async = require('async');

module.exports = HomepageController;
function HomepageController(properties) {
  this._talks = properties.talks;
}

HomepageController.prototype.index = function(req, res) {
  async.parallel({
    recentTalks: this._talks.findRecent,
  }, function(err, results) {
    res.render('homepage/index', _.extend(results, {
      title: 'Felix Geisendörfer',
    }));
  });
};
