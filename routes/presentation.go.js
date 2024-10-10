const NUMBERS_ONLY = /^[0-9]*$/;
function inRange(x, min, max) {
  return (x - min) * (x - max) <= 0;
}
