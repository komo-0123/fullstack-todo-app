import Footer from "@/components/layout/Footer";
import Header from "@/components/layout/Header";
import { ErrorModal, ErrorModalProvider } from "@/features/ErrorModal";
import { TodoInput, TodoList } from "@/features/Todo";

function App() {
  return (
    <ErrorModalProvider>
      <div className="max-w-3xl mx-auto p-6 bg-white mt-12">
        <Header />
        <main>
          <TodoInput />
          <TodoList />
        </main>
        <Footer />
        <ErrorModal />
      </div>
    </ErrorModalProvider>
  );
}

export default App;
