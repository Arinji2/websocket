import { Toaster } from "react-hot-toast";
import { Route, Routes } from "react-router-dom";
import Home from "./Home";
function App() {
  return (
    <>
      <div>
        <Toaster
          position="top-right"
          toastOptions={{
            style: {
              borderRadius: "10px",
              background: "#333",
              color: "#fff",
            },
          }}
        />
      </div>
      <div className="w-full min-h-[100svh] bg-slate-900 flex flex-col items-center justify-center">
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </div>
    </>
  );
}

export default App;
