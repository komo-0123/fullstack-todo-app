import useSWR from "swr";
import { TodosResponse } from "@/types";
import { API } from "@/constant";
import TodoItem from "./TodoItem";

const TodoList = () => {
  // useSWRでデータを取得する
  const endPoint = `${API.BASE_URL}${API.TODOS}`;
  const fetcher = (url: string) => fetch(url).then((res) => res.json());
  const { data, error, isLoading } = useSWR<TodosResponse>(endPoint, fetcher);

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div>Error...</div>;
  if (!data) return <div>No data...</div>;

  return (
    <ul className="space-y-2">
      {data.data.map((todo) => (
        <TodoItem key={todo.id} todo={todo} />
      ))}
    </ul>
  );
};

export default TodoList;
