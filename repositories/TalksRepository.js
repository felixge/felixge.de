var talks = require('../data/talks');

module.exports = TalksRepository;
function TalksRepository() {
}

TalksRepository.prototype.findUpcoming = function(cb) {
  cb(null, talks);
};

TalksRepository.prototype.findPast = function(cb) {
  cb(null, talks);
};
