import { useContext, useState } from "react";
import { TodosResponse } from "@/types";
import { API } from "@/constant";
import { mutate } from "swr";
import { ErrorModalContext } from "@/features/ErrorModal";

const TodoInput = () => {
  const { showError } = useContext(ErrorModalContext);
  const [inputValue, seInputValue] = useState("");

  // jsonでリクエスを送信。サーバーはlocalhost:8000/todosにリクエストを受け取る
  const addTodo = async (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();

    const endPoint = `${API.BASE_URL}${API.TODOS}`;
    const res = await fetch(endPoint, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        title: inputValue,
        is_complete: false,
      }),
    });

    const data: TodosResponse = await res.json();

    if (data.status.error) {
      showError(data.status.error_message);
    }

    mutate(`${API.BASE_URL}${API.TODOS}`);
    seInputValue("");
  };

  return (
    <form className="flex gap-2 mb-4" onSubmit={(e) => addTodo(e)}>
      <input
        type="text"
        value={inputValue}
        onChange={(e) => seInputValue(e.target.value)}
        className="flex-1 p-2 border border-teal-400 rounded focus:outline-none focus:ring-2 focus:ring-teal-500"
        placeholder="新しいTODOを追加"
      />
      <button className="px-4 py-2 bg-teal-500 text-white font-semibold rounded hover:bg-teal-600 transition">
        追加
      </button>
    </form>
  );
};

export default TodoInput;
