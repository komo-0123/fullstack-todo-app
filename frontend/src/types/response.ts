type Data = {
  id: number;
  title: string;
  is_complete: boolean;
};

type TodoResponse = {
  data: Data | null;
  status: {
    code: number;
    error: boolean;
    error_message: string;
  };
};

type TodosResponse = {
  data: Data[];
  status: {
    code: number;
    error: boolean;
    error_message: string;
  };
};

export type { Data, TodoResponse, TodosResponse };
