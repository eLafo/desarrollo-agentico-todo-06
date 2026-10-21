// Pide confirmación antes de enviar formularios marcados con data-confirm.
document.addEventListener("submit", function (event) {
  var form = event.target;
  if (form.dataset.confirm && !window.confirm(form.dataset.confirm)) {
    event.preventDefault();
  }

  if (!form.matches("[data-task-toggle]")) {
    return;
  }

  event.preventDefault();

  var task = form.closest(".task");
  var completed = task.classList.toggle("task-completed");
  var button = form.querySelector("button");
  var status = task.querySelector(".task-status");

  var label = completed ? "Reabrir" : "Completar";
  button.querySelector(".toggle-label").textContent = label;
  button.setAttribute("aria-pressed", String(completed));
  button.setAttribute("aria-label", label + " tarea: " + task.querySelector(".task-title").textContent);
  button.setAttribute("title", label + " tarea");

  if (completed && !status) {
    status = document.createElement("span");
    status.className = "task-status";
    status.textContent = "Completada";
    task.insertBefore(status, task.querySelector(".task-actions"));
  } else if (!completed && status) {
    status.remove();
  }

  fetch(form.action);
});

// Ordena la lista arrastrando cada tarea por su asa, como mejora sobre los botones «Subir» y
// «Bajar». Al soltar, envía el formulario de mover de la tarea con los puestos que se ha
// desplazado y pinta la lista que devuelve el servidor, que es quien decide si se puede mover.
var dragging = null;

function enableSorting(list) {
  list.classList.add("task-list-sortable");
  list.querySelectorAll("[data-task-handle]").forEach(function (handle) {
    handle.draggable = true;
  });
}

function taskIndex(task) {
  return Array.prototype.indexOf.call(task.parentNode.children, task);
}

function insideDraggedList(target) {
  return dragging && target instanceof Element && dragging.list.contains(target);
}

function saveMove(list, task, offset) {
  var form = task.querySelector("[data-task-move]");
  var body = new URLSearchParams(new FormData(form));
  body.set("offset", String(offset));

  fetch(form.action, { method: "POST", body: body })
    .then(function (response) {
      if (!response.ok) {
        throw new Error(response.statusText);
      }
      return response.text();
    })
    .then(function (html) {
      var saved = new DOMParser().parseFromString(html, "text/html").querySelector(".task-list");
      if (!saved) {
        throw new Error("La respuesta no trae la lista.");
      }
      list.replaceWith(saved);
      enableSorting(saved);
    })
    .catch(function () {
      window.location.reload();
    });
}

document.addEventListener("dragstart", function (event) {
  var handle = event.target instanceof Element && event.target.closest("[data-task-handle]");
  if (!handle) {
    return;
  }

  var task = handle.closest(".task");
  var box = task.getBoundingClientRect();
  var current = {
    task: task,
    list: task.parentNode,
    from: taskIndex(task),
    next: task.nextElementSibling,
    dropped: false
  };
  dragging = current;

  event.dataTransfer.effectAllowed = "move";
  event.dataTransfer.setData("text/plain", "");
  event.dataTransfer.setDragImage(task, event.clientX - box.left, event.clientY - box.top);

  window.requestAnimationFrame(function () {
    if (dragging === current) {
      task.classList.add("task-dragging");
    }
  });
});

document.addEventListener("dragover", function (event) {
  if (!insideDraggedList(event.target)) {
    return;
  }

  event.preventDefault();
  event.dataTransfer.dropEffect = "move";

  var over = event.target.closest(".task");
  if (!over || over === dragging.task || over.parentNode !== dragging.list) {
    return;
  }

  var box = over.getBoundingClientRect();
  var before = event.clientY < box.top + box.height / 2;
  dragging.list.insertBefore(dragging.task, before ? over : over.nextElementSibling);
});

document.addEventListener("drop", function (event) {
  if (!insideDraggedList(event.target)) {
    return;
  }

  event.preventDefault();
  dragging.dropped = true;
});

document.addEventListener("dragend", function () {
  if (!dragging) {
    return;
  }

  var drag = dragging;
  dragging = null;
  drag.task.classList.remove("task-dragging");

  if (!drag.dropped) {
    drag.list.insertBefore(drag.task, drag.next);
    return;
  }

  var offset = taskIndex(drag.task) - drag.from;
  if (offset !== 0) {
    saveMove(drag.list, drag.task, offset);
  }
});

var taskList = document.querySelector(".task-list");
if (taskList) {
  enableSorting(taskList);
}
