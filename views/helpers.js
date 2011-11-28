var jade = require('jade');

exports.days = function(date) {
  if (!date) return date;

  var from = date;
  var to   = date;

  if (Array.isArray(date)) {
    from = date[0];
    to   = date[1];
  }

  from = this._dateToJson(from);
  to   = this._dateToJson(to);

  if (from.year !== to.year) {
    return(
      from.month + ' ' + from.day + ', ' + from.year + ' - ' +
      to.month + ' ' + to.day + ', ' + to.year
    );
  }

  if (from.month !== to.month) {
    return(
      from.month + ' ' + from.day + ' - ' +
      to.month + ' ' + to.day + ', ' + to.year
    );
  }

  if (from.day !== to.day) {
    return(
      from.month + ' ' + from.day + ' - ' +
      to.day + ', ' + to.year
    );
  }

  return from.month + ' ' + from.day + ', ' + from.year;
};

exports._dateToJson = function(date) {
  var months = [
    'January',
    'Feburary',
    'March',
    'April',
    'May',
    'June',
    'July',
    'August',
    'September',
    'October',
    'November',
    'December',
  ];

  var month = months[date.getMonth()];
  var day   = date.getDate();
  var year  = date.getFullYear();

  return {
    day   : day,
    month : month,
    year  : year,
  };
};

exports.link = function(title, url) {
  if (!url) return title;

  // @todo escape title

  return '<a href="' + url + '">' + title + '</a>';
};
