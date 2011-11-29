var talks = require('../data/talks');

module.exports = TalksRepository;
function TalksRepository() {
}

TalksRepository.prototype.findUpcoming = function(cb) {
  cb(null, talks.filter(function(talk) {
    return !TalksRepository._isPast(talk);
  }));
};

TalksRepository.prototype.findRecent = function(cb) {
  cb(null, talks.slice(0, 5));
};

TalksRepository.prototype.findTalks = function(cb) {
  cb(null, talks);
};

TalksRepository._isPast = function(talk) {
  var date = talk.date[1] || talk.date;
  return Date.now() > date;
};
