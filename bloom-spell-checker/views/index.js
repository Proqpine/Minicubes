const nameInput = document.getElementById("name");
const result = document.getElementById("result");

nameInput.addEventListener("input", function () {
  const word = nameInput.value;

  fetch("/submit", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ word: word }),
  })
    .then((response) => response.text())
    .then((data) => {
      result.innerHTML = data;
    })
    .catch((error) => console.error("Error:", error));
});
