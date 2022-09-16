// utility function for parsing the error message
export const parseErrorMessage = (message: string, key: string) => {
  const firstWord = message.split(" ")[0];
  const firstTwoWords = message.split(" ").slice(0, 2).join(" ");
  if (firstWord.toLowerCase() === "value") {
    return message.replace("Value", key);
  }
  if (firstTwoWords.toLowerCase() === "this field") {
    return message.replace("This field", key);
  }
  return message;
};
