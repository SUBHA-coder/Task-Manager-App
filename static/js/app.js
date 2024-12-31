const apiUrl = "/tasks";

async function fetchTasks() {
  const response = await fetch(apiUrl);
  const tasks = await response.json();
  renderTasks(tasks);
}

function renderTasks(tasks) {
  const taskList = document.getElementById("taskList");
  taskList.innerHTML = "";

  tasks.forEach(task => {
    const taskElement = document.createElement("div");
    taskElement.className = "flex justify-between items-center p-4 border-b border-gray-200";

    const taskInfo = `
      <div>
        <h3 class="text-lg font-bold">${task.title}</h3>
        <p class="text-gray-500">Priority: ${task.priority} | Due: ${task.dueDate}</p>
        <p class="text-sm text-gray-500">Status: ${task.complete ? "Complete" : "Incomplete"}</p>
      </div>
    `;

    const taskActions = `
      <div class="flex space-x-2">
        <button onclick="toggleComplete(${task.id}, ${!task.complete})" class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded">${task.complete ? "Mark Incomplete" : "Mark Complete"}</button>
        <button onclick="deleteTask(${task.id})" class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded">Delete</button>
      </div>
    `;

    taskElement.innerHTML = taskInfo + taskActions;
    taskList.appendChild(taskElement);
  });
}

async function addTask(event) {
    event.preventDefault();
    const title = document.getElementById("taskTitle").value;
    const priority = document.getElementById("taskPriority").value;
    const dueDate = document.getElementById("taskDueDate").value;
  
    const newTask = { title, priority, dueDate };
  
    try {
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(newTask),
      });
  
      if (!response.ok) {
        const errorText = await response.text();
        console.error("Error adding task:", errorText);
        alert("Failed to add task: " + errorText);
        return;
      }
  
      // Fetch and refresh the task list
      fetchTasks();
      document.getElementById("addTaskForm").reset();
    } catch (error) {
      console.error("Network or server error:", error);
      alert("An unexpected error occurred. Check the console for details.");
    }
  }
  

async function toggleComplete(taskId, newStatus) {
  const response = await fetch(`${apiUrl}/${taskId}`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ complete: newStatus })
  });

  if (response.ok) {
    fetchTasks();
  } else {
    alert("Failed to update task!");
  }
}

async function deleteTask(taskId) {
  const response = await fetch(`${apiUrl}/${taskId}`, { method: "DELETE" });

  if (response.ok) {
    fetchTasks();
  } else {
    alert("Failed to delete task!");
  }
}

document.getElementById("addTaskForm").addEventListener("submit", addTask);

// Fetch tasks on page load
fetchTasks();
