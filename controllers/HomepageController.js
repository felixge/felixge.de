module.exports = HomepageController;
function HomepageController(properties) {
  this.talks = properties.talks;
}

HomepageController.prototype.index = function(req, res) {
  res.render('homepage', {
    title: 'Felix Geisendörfer',
    talks: {}
  });
};
