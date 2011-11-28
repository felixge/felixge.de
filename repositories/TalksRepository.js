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
  cb(null, talks);
};

TalksRepository.prototype.findPast = function(cb) {
  cb(null, talks.filter(function(talk) {
    return TalksRepository._isPast(talk);
  }));
};

TalksRepository._isPast = function(talk) {
  var date = talk.date[1] || talk.date;
  return Date.now() > date;
};
