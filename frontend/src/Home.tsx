import { UserSchema } from "./types/user";

export default function Home() {
  return (
    <div className="w-full h-full flex flex-col items-center gap-4 justify-center">
      <h1 className="text-6xl font-bold text-white">CREATE USER</h1>
      <form
        onSubmit={async (e) => {
          e.preventDefault();
          const formData = new FormData(e.currentTarget);
          const name = formData.get("name");
          const email = formData.get("email");
          if (!email || !name) {
            alert("Please fill in all fields");
            return;
          }

          const backendDomain = import.meta.env.VITE_BACKEND_URL;
          const response = await fetch(`${backendDomain}/api/user/create`, {
            method: "POST",
            headers: {
              "Content-Type": "application/json",
            },
            body: JSON.stringify({
              name: name,
              email: email,
            }),
          });

          if (response.ok) {
            alert("User created successfully");
          } else {
            alert("Error creating user");
            console.error(response);
          }
          try {
            const rawUser = await response.json();
            const parsedUser = UserSchema.parse(rawUser);
            localStorage.setItem("user", JSON.stringify(parsedUser));
          } catch (error) {
            console.error(error);
          }
        }}
        className="flex flex-col items-center justify-center gap-3"
      >
        <input
          type="text"
          name="name"
          placeholder="name..."
          className="w-full p-2 text-white border-b-2 border-black shadow-black bg-transparent"
        />
        <input
          type="email"
          name="email"
          placeholder="Email..."
          className="w-full p-2 text-white border-b-2 border-black shadow-black bg-transparent"
        />
        <button
          type="submit"
          className="bg-blue-500 text-white px-4 py-2 rounded-md shadow-md shadow-black font-bold"
        >
          Create User
        </button>
      </form>
    </div>
  );
}
