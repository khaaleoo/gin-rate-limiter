const validTypes = ["build", "chore", "ci", "docs", "feat", "fix", "perf", "refactor", "revert", "style", "test"];
const strEmojiPattern = "^(:[a-z]+_[a-z]+:|\\p{Emoji}) ";
const reEmoji = RegExp(strEmojiPattern, "u");
const isValidType = (header) => {
  const isValid = header.length > 0 && reEmoji.test(header);
  if (isValid) {
    return [true];
  }
  return [false, "type must be present"];
};

const isValidSubject = (header) => {
  let re = RegExp(`${strEmojiPattern}(${validTypes.join`|`})(:|\\()`, "u");
  const isValid = header.length > 0 && re.test(header);
  if (isValid) {
    return [true];
  }
  return [false, "subject must be present"];
};

module.exports = {
  extends: ["@commitlint/config-conventional"],
  plugins: ["commitlint-plugin-function-rules"],
  rules: {
    "body-max-line-length": [1, "always", Infinity],
    "header-max-length": [1, "always", Infinity],
    "subject-empty": [0],
    "type-empty": [0],
    "function-rules/type-empty": [
      2, // level: error
      "always",
      (parsed) => isValidType(parsed.header),
    ],
    "function-rules/subject-empty": [
      2, // level: error
      "always",
      (parsed) => isValidSubject(parsed.header),
    ],
  },
};
