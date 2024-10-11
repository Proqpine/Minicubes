// Parts of this script was gotten from
// https://webdesign.tutsplus.com/learn-how-to-code-a-simple-javascript-calendar-and-datepicker--cms-108322t

const socket = new WebSocket("ws://localhost:8080/ws");

function requestAvailabilty(dateTo) {
  socket.send(
    JSON.stringify({
      type: "REQUEST_AVAILABILITY",
      payload: {
        date: dateTo, // format: "2024-10-11"
      },
    }),
  );
}

socket.onmessage = function (event) {
  const message = JSON.parse(event.data);
  if (message.type === "AVAILABILITY_RESPONSE") {
    const availableTimes = message.payload.availableTimes;
    console.log("Available times:", availableTimes);
  }
};

socket.onclose = (event) => {
  console.log("WebSocket connection closed");
};

socket.onerror = (error) => {
  console.error("WebSocket error:", error);
};

let display = document.querySelector(".display");
let previous = document.querySelector(".left");
let next = document.querySelector(".right");
let days = document.querySelector(".days");
let selected = document.querySelector(".selected");
let day = document.querySelector(".day");

let date = new Date();
let year = date.getFullYear();
let month = date.getMonth();

function displayCalendar() {
  const firstDay = new Date(year, month, 0);
  const firstDayIndex = firstDay.getDay();

  const lastDay = new Date(year, month + 1, 0);
  const numberOfDays = lastDay.getDate();

  let formattedDate = date.toLocaleString("en-UK", {
    month: "long",
    year: "numeric",
  });
  display.innerHTML = `${formattedDate}`;

  for (let x = 1; x <= firstDayIndex; x++) {
    let emptyDay = document.createElement("div");
    emptyDay.className = "bg-white";
    let button = document.createElement("button");
    button.type = "button";
    button.className =
      "mx-auto flex size-10 w-full items-center justify-center text-gray-400 hover:text-blue-600";
    let time = document.createElement("time");
    button.appendChild(time);
    emptyDay.appendChild(button);
    days.appendChild(emptyDay);
  }

  for (let i = 1; i <= numberOfDays; i++) {
    let currentDate = new Date(year, month, i);

    let dayDiv = document.createElement("div");
    dayDiv.className = "bg-white";

    const button = document.createElement("button");
    button.type = "button";
    button.className =
      "mx-auto flex size-10 w-full items-center justify-center text-gray-400 hover:text-blue-600";

    const time = document.createElement("time");
    time.datetime = currentDate.toISOString().split("T")[0];
    time.textContent = i;

    button.appendChild(time);
    dayDiv.appendChild(button);
    days.appendChild(dayDiv);

    button.dataset.date = currentDate.toDateString();
    days.dataset.date = currentDate.toDateString();
    time.dataset.date = currentDate.toDateString();

    if (
      currentDate.getFullYear() === new Date().getFullYear() &&
      currentDate.getMonth() === new Date().getMonth() &&
      currentDate.getDate() === new Date().getDate()
    ) {
      dayDiv.classList.add("current-date");
    }
  }
}
displayCalendar();

function displaySelected() {
  const dayElements = document.querySelectorAll(".days button");
  dayElements.forEach((day) => {
    day.addEventListener("click", (e) => {
      dayElements.forEach((d) => {
        d.classList.remove("bg-blue-600");
        d.classList.remove("text-white");
        d.classList.add("text-gray-400");
        d.classList.add("hover:text-blue-600");
      });
      const selectedButton = e.target.closest("button");
      selectedButton.classList.remove("text-gray-400");
      selectedButton.classList.remove("hover:text-blue-600");
      selectedButton.classList.add("bg-blue-600");
      selectedButton.classList.add("text-white");

      const selectedDate = e.target.dataset.date;
      const dateObj = new Date(selectedDate);

      // Format the date as "YYYY-MM-DD" in local time
      const year = dateObj.getFullYear();
      const month = String(dateObj.getMonth() + 1).padStart(2, "0"); // Months are 0-based
      const day = String(dateObj.getDate()).padStart(2, "0");

      const formattedDate = `${year}-${month}-${day}`;
      console.log(formattedDate);
      requestAvailabilty(formattedDate);
    });
  });
}
displaySelected();

previous.addEventListener("click", () => {
  days.innerHTML = "";
  selected.innerHTML = "";

  if (month < 0) {
    month = 11;
    year = year - 1;
  }

  month = month - 1;
  date.setMonth(month);

  displayCalendar();
  displaySelected();
});

next.addEventListener("click", () => {
  days.innerHTML = "";
  selected.innerHTML = "";
  if (month > 11) {
    month = 0;
    year = year + 1;
  }
  month = month + 1;
  date.setMonth(month);
  displayCalendar();
  displaySelected();
});
