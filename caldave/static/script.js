// Parts of this script was gotten from
// https://webdesign.tutsplus.com/learn-how-to-code-a-simple-javascript-calendar-and-datepicker--cms-108322t

const socket = new WebSocket("ws://localhost:8080/ws");

function requestAvailability(dateTo) {
  sendMessage({
    type: "REQUEST_AVAILABILITY",
    payload: {
      date: dateTo,
    },
  });
}

function updateEvents(startDate, endDate) {
  sendMessage({
    type: "UPDATE_AVAILABILITY",
    payload: {
      startDate: startDate,
      endDate: endDate,
    },
  });
}

let isWebSocketReady = false;
const pendingMessages = [];

socket.onopen = (event) => {
  console.log("WebSocket connection established");
  isWebSocketReady = true;
  sendPendingMessages();
  displayCalendar();
  displaySelected();
};

function sendMessage(message) {
  if (isWebSocketReady) {
    socket.send(JSON.stringify(message));
  } else {
    pendingMessages.push(message);
  }
}

function sendPendingMessages() {
  while (pendingMessages.length > 0) {
    const message = pendingMessages.shift();
    socket.send(JSON.stringify(message));
  }
}

socket.onmessage = (event) => {
  const message = JSON.parse(event.data);
  if (message.type === "AVAILABILITY_RESPONSE") {
    const availableTimes = message.payload.availableTimes;
    console.log("Available times:", availableTimes);
    displayAvailableTimes(availableTimes);
  } else if (message.type === "EVENTS_UPDATED") {
    console.log("Events updated successfully");
  }
};

socket.onclose = (event) => {
  console.log("WebSocket connection closed");
  isWebSocketReady = false;
};

socket.onerror = (error) => {
  console.error("WebSocket error:", error);
  isWebSocketReady = false;
};

const display = document.querySelector(".display");
const previous = document.querySelector(".left");
const next = document.querySelector(".right");
const days = document.querySelector(".days");
const selected = document.querySelector(".selected");
const day = document.querySelector(".day");

const date = new Date();
let year = date.getFullYear();
let month = date.getMonth();

function displayCalendar() {
  const firstDay = new Date(year, month, 0);
  const firstDayIndex = firstDay.getDay();
  const lastDay = new Date(year, month + 1, 0);
  const numberOfDays = lastDay.getDate();

  const formattedDate = date.toLocaleString("en-UK", {
    month: "long",
    year: "numeric",
  });
  display.innerHTML = `${formattedDate}`;

  for (let x = 1; x <= firstDayIndex; x++) {
    const emptyDay = document.createElement("div");
    emptyDay.className = "bg-white";
    const button = document.createElement("button");
    button.type = "button";
    button.className =
      "mx-auto flex size-10 w-full items-center justify-center text-gray-400 hover:text-blue-600";
    const time = document.createElement("time");
    button.appendChild(time);
    emptyDay.appendChild(button);
    days.appendChild(emptyDay);
  }

  for (let i = 1; i <= numberOfDays; i++) {
    const currentDate = new Date(year, month, i);

    const dayDiv = document.createElement("div");
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
  updateEventsForCurrentMonth();
}

function updateEventsForCurrentMonth() {
  const startDate = new Date(year, month, 1);
  const endDate = new Date(year, month + 1, 0);
  updateEvents(
    startDate.toISOString().split("T")[0],
    endDate.toISOString().split("T")[0],
  );
}

function displaySelected() {
  const dayElements = document.querySelectorAll(".days button");
  const selectedDay = document.querySelector(".selected-date");
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
      selectedDay.innerHTML = `${selectedDate}`;
      const dateObj = new Date(selectedDate);

      const year = dateObj.getFullYear();
      const month = String(dateObj.getMonth() + 1).padStart(2, "0");
      const day = String(dateObj.getDate()).padStart(2, "0");

      const formattedDate = `${year}-${month}-${day}`;
      console.log(formattedDate);
      requestAvailability(formattedDate);
    });
  });
}

function displayAvailableTimes(availableTimes) {
  const timeSlotsContainer = document.querySelector(".time-slots");
  timeSlotsContainer.innerHTML = "";

  availableTimes.forEach((timeSlot) => {
    const start = timeSlot.start;
    const end = timeSlot.end;
    const startTime = new Date(`2000-01-01T${start}`);
    const endTime = new Date(`2000-01-01T${end}`);

    while (startTime < endTime) {
      const slotStart = startTime.toTimeString().slice(0, 5);
      startTime.setMinutes(startTime.getMinutes() + 30);
      const slotEnd = startTime.toTimeString().slice(0, 5);

      const timeSlotDiv = document.createElement("div");
      // class="available text-gray-500 hover:text-gray-800 cursor-pointer border-[1.5px] border-gray-400 px-2 text-center py-1 rounded-lg"
      timeSlotDiv.className = "flex items-center justify-center";
      timeSlotDiv.innerHTML = `
           <span class="text-gray-500 w-full inline-block hover:text-gray-800 cursor-pointer border-[1.5px] border-gray-400 px-2 text-center py-1 rounded-lg transition-colors duration-200 ease-in-out hover:bg-gray-100">
             ${slotStart} - ${slotEnd}
           </span>
         `;
      timeSlotsContainer.appendChild(timeSlotDiv);
    }
  });
}

document.addEventListener("DOMContentLoaded", () => {
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
});
