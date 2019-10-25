export enum Platform {
  Unknown,
  Mac,
  Ios,
  Windows,
  Android,
  Linux,
}

// adapted from
// https://stackoverflow.com/questions/38241480/detect-macos-ios-windows-android-and-linux-os-with-js
export function getPlatform() {
  const userAgent = window.navigator.userAgent;
  const platform = window.navigator.platform;
  const macosPlatforms = ['Macintosh', 'MacIntel', 'MacPPC', 'Mac68K'];
  const windowsPlatforms = ['Win32', 'Win64', 'Windows', 'WinCE'];
  const iosPlatforms = ['iPhone', 'iPad', 'iPod'];
  if (macosPlatforms.indexOf(platform) !== -1) {
    return Platform.Mac;
  } else if (iosPlatforms.indexOf(platform) !== -1) {
    return Platform.Ios;
  } else if (windowsPlatforms.indexOf(platform) !== -1) {
    return Platform.Windows;
  } else if (/Android/.test(userAgent)) {
    return Platform.Android;
  } else if (/Linux/.test(platform)) {
    return Platform.Linux;
  }
  return Platform.Unknown;
}
