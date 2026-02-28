import React, { useEffect, useState, useContext } from "react";
import axios from "axios";
import { AuthContext } from "./AuthContext";

interface Todo {
  id: number;
  title: string;
  completed: boolean;
}

const Todos: React.FC = () => {
  const auth = useContext(AuthContext);

  const [todos, setTodos] = useState<Todo[]>([]);
  const [newTodo, setNewTodo] = useState("");
  const [error, setError] = useState("");

  // Inline edit state
  const [editTodoId, setEditTodoId] = useState<number | null>(null);
  const [editTitle, setEditTitle] = useState("");

  // Filter state
  const [filter, setFilter] = useState<"All" | "Active" | "Completed">("All");

  const config = {
    headers: {
      Authorization: `Bearer ${auth?.token}`,
    },
  };

  // 🔹 Fetch Todos
  const fetchTodos = async () => {
    if (!auth?.token) return;
    try {
      const res = await axios.get("http://localhost:8080/todos", config);
      setTodos(res.data);
    } catch (err: any) {
      console.error(err);
      setError("Failed to fetch todos");
    }
  };

  useEffect(() => {
    fetchTodos();
  }, []);

  // 🔹 Add Todo
  const addTodo = async () => {
    if (!newTodo.trim()) return;

    try {
      const res = await axios.post(
        "http://localhost:8080/todos",
        { title: newTodo, completed: false },
        config
      );
      setTodos([...todos, res.data]);
      setNewTodo("");
    } catch (err: any) {
      console.error(err);
      setError("Failed to add todo");
    }
  };

  // 🔹 Toggle Complete
  const toggleTodo = async (todo: Todo) => {
    const oldTodos = [...todos];
    setTodos(
      todos.map((t) =>
        t.id === todo.id ? { ...t, completed: !t.completed } : t
      )
    );

    try {
      await axios.put(
        `http://localhost:8080/todos/${todo.id}`,
        { ...todo, completed: !todo.completed },
        config
      );
    } catch (err: any) {
      console.error(err);
      setTodos(oldTodos); // rollback
      setError("Failed to update todo");
    }
  };

  // 🔹 Delete Todo
  const deleteTodo = async (todo: Todo) => {
    if (!window.confirm("Really delete this todo?")) return;

    const oldTodos = [...todos];
    setTodos(todos.filter((t) => t.id !== todo.id));

    try {
      await axios.delete(
        `http://localhost:8080/todos/${todo.id}`,
        config
      );
    } catch (err: any) {
      console.error(err);
      setError("Failed to delete todo");
      setTodos(oldTodos); // rollback
    }
  };

  // 🔹 Start Edit
  const startEdit = (todo: Todo) => {
    setEditTodoId(todo.id);
    setEditTitle(todo.title);
  };

  // 🔹 Save Edit
  const saveEdit = async (todo: Todo) => {
    if (!editTitle.trim()) return;

    const oldTodos = [...todos];

    setTodos(
      todos.map((t) =>
        t.id === todo.id ? { ...t, title: editTitle } : t
      )
    );
    setEditTodoId(null);

    try {
      await axios.put(
        `http://localhost:8080/todos/${todo.id}`,
        { ...todo, title: editTitle },
        config
      );
    } catch (err: any) {
      console.error(err);
      setTodos(oldTodos);
      setError("Failed to update todo");
    }
  };

  // 🔹 Filtered Todos
  const filteredTodos = todos.filter((todo) => {
    if (filter === "Active") return !todo.completed;
    if (filter === "Completed") return todo.completed;
    return true;
  });

  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center p-6">
      <h1 className="text-3xl font-bold mb-6">My Todos 📝</h1>

      {error && (
        <div
          className="bg-red-100 text-red-600 px-4 py-2 rounded mb-4 cursor-pointer"
          onClick={() => setError("")}
        >
          {error}
        </div>
      )}

      {/* Add Todo */}
      <div className="flex gap-2 mb-6 w-full max-w-md">
        <input
          type="text"
          placeholder="Add new todo..."
          value={newTodo}
          onChange={(e) => setNewTodo(e.target.value)}
          onKeyDown={(e) => e.key === "Enter" && addTodo()}
          className="flex-1 border px-3 py-2 rounded"
        />
        <button
          onClick={addTodo}
          className="bg-indigo-600 text-white px-4 py-2 rounded hover:bg-indigo-700"
        >
          Add
        </button>
      </div>

      {/* Filter Buttons */}
      <div className="mb-4 flex gap-2">
        {["All", "Active", "Completed"].map((f) => (
          <button
            key={f}
            className={`px-3 py-1 rounded ${
              filter === f
                ? "bg-indigo-600 text-white"
                : "bg-gray-200"
            }`}
            onClick={() => setFilter(f as any)}
          >
            {f}
          </button>
        ))}
      </div>

      {/* Todo List */}
      <ul className="w-full max-w-md max-h-[400px] overflow-y-auto bg-white shadow rounded-lg p-4 space-y-3">
        {filteredTodos.map((todo) => (
          <li
            key={todo.id}
            className="flex justify-between items-center border-b pb-2"
          >
            {editTodoId === todo.id ? (
              <input
                value={editTitle}
                onChange={(e) => setEditTitle(e.target.value)}
                onBlur={() => saveEdit(todo)}
                onKeyDown={(e) =>
                  e.key === "Enter" && saveEdit(todo)
                }
                className="border px-2 py-1 rounded w-full mr-2"
                autoFocus
              />
            ) : (
              <span
                onDoubleClick={() => startEdit(todo)}
                onClick={() => toggleTodo(todo)}
                className={`cursor-pointer transition-all duration-300 ${
                  todo.completed
                    ? "line-through text-gray-400"
                    : "text-gray-800"
                }`}
              >
                {todo.title}
              </span>
            )}

            <button
              onClick={() => deleteTodo(todo)}
              className="text-red-500 hover:text-red-700 ml-3"
            >
              ✕
            </button>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Todos;