import { useNavigate } from "react-router-dom";

export default function Home() {
  const navigate = useNavigate();

  return (
    <div className="home">
      <p>home</p>
      <button onClick={() => navigate("/upload")} role="button" tabIndex={-1}>
        upload
      </button>
      <button onClick={() => navigate("/download")} role="button" tabIndex={-1}>
        download
      </button>
    </div>
  );
}
