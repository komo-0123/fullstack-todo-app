import { useEffect, useRef, useState } from "react";
import { mutate } from "swr";
import { API } from "@/constant";
import { Data } from "@/types";

const TodoItem = ({ todo }: { todo: Data }) => {
  const [isEditing, setIsEditing] = useState(false);
  const [title, setTitle] = useState(todo.title);
  const [isComplete, setIsComplete] = useState(todo.is_complete);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (isEditing) {
      inputRef.current?.focus();
    }
  }, [isEditing]);

  // TODO更新APIを呼び出し、ミューテートする
  const updateTodo = async (newTitle: string, newIsComplete: boolean) => {
    // 両方の値が変更されていない場合は何もしない
    if (newTitle === todo.title && newIsComplete === todo.is_complete) {
      return;
    }

    const endPoint = `${API.BASE_URL}${API.TODOS}/${todo.id}`;
    await fetch(endPoint, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ title: newTitle, is_complete: newIsComplete }),
    });

    await mutate(`${API.BASE_URL}${API.TODOS}`);
  };

  // チェックボックスの状態を更新する
  const handleIsCompleteUpdate = async (e: React.ChangeEvent<HTMLInputElement>) => {
    setIsComplete(e.target.checked);
    await updateTodo(title, e.target.checked);
  };

  // タイトルの更新を確定する
  const handleTitleUpdate = async () => {
    setIsEditing(false);
    await updateTodo(title, isComplete);
  };

  // TODOを削除する
  const handleDelete = async () => {
    const endPoint = `${API.BASE_URL}${API.TODOS}/${todo.id}`;
    await fetch(endPoint, {
      method: "DELETE",
    });

    await mutate(`${API.BASE_URL}${API.TODOS}`);
  };

  return (
    <li key={todo.id} className="flex gap-x-2 items-center p-2 border-b border-b-teal-600">
      <div className="flex flex-1 items-center gap-x-2">
        <input
          id={`todo${todo.id}`}
          type="checkbox"
          checked={isComplete}
          onChange={(e) => handleIsCompleteUpdate(e)}
          className="w-5 h-5 text-teal-500 rounded focus:ring-teal-500"
        />
        {isEditing ? (
          <input
            type="text"
            value={title}
            ref={inputRef}
            onChange={(e) => setTitle(e.target.value)}
            onBlur={handleTitleUpdate}
            className="w-full px-2 py-1 border border-gray-300 rounded focus:outline-none focus:ring focus:ring-teal-500"
          />
        ) : (
          <label htmlFor={`todo${todo.id}`} className="text-gray-700">
            {todo.title}
          </label>
        )}
      </div>
      <div>
        <button
          onClick={() => setIsEditing(true)}
          className="px-3 py-1 mr-1 text-sm text-white bg-blue-800 rounded hover:bg-blue-900 transition"
        >
          編集
        </button>
        <button
          onClick={handleDelete}
          className="px-3 py-1 text-sm text-white bg-red-500 rounded hover:bg-red-600 transition"
        >
          削除
        </button>
      </div>
    </li>
  );
};

export default TodoItem;
