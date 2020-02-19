// adapted from:
// https://stackoverflow.com/questions/5968196/check-cookie-if-cookie-exists
export function getCookie(name: string): string | undefined {
  const dc = document.cookie;
  const prefix = name + '=';
  let begin = dc.indexOf('; ' + prefix);
  let end = undefined;
  if (begin == -1) {
    begin = dc.indexOf(prefix);
    if (begin != 0) {
      return undefined;
    }
  } else {
    begin += 2;
    end = document.cookie.indexOf(';', begin);
    if (end == -1) {
      end = dc.length;
    }
  }
  // because unescape has been deprecated, replaced with decodeURI
  // return unescape(dc.substring(begin + prefix.length, end));
  return decodeURI(dc.substring(begin + prefix.length, end));
}
