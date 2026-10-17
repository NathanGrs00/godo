import React, { useState } from "react";
import { useNavigate } from "react-router-dom";
import SearchIcon from "../../assets/icons/NavBar/SearchIcon.tsx";
import PageLayout from "../Common/PageLayout.tsx";
import { useData } from "../../context/dataContext.ts";
import { search } from "../../utils/search.ts";
import { relativeTime } from "../../utils/date.ts";
import "../../styles/Pages/Pages.css";

/** Searches tasks and deadlines by title, description and tag name. */
const SearchPage: React.FC = () => {
    const data = useData();
    const navigate = useNavigate();
    const [query, setQuery] = useState("");
    const results = search(query, data);
    const nothing = query.trim() !== "" && results.tasks.length === 0 && results.deadlines.length === 0;

    return (
        <PageLayout title="Search" icon={<SearchIcon />}>
            <input
                type="search"
                className="search-input"
                aria-label="Search"
                placeholder="Search tasks and deadlines..."
                value={query}
                onChange={e => setQuery(e.target.value)}
                autoFocus
            />
            {nothing && <p className="empty-message">Nothing found for "{query}".</p>}

            {results.tasks.length > 0 && (
                <section className="page-section">
                    <h2>Tasks</h2>
                    <ul className="result-list">
                        {results.tasks.map(task => (
                            <li key={task.id}>
                                <button type="button" className="link-button" onClick={() => navigate("/")}>
                                    {task.completed ? "✓ " : ""}{task.title}
                                </button>
                            </li>
                        ))}
                    </ul>
                </section>
            )}

            {results.deadlines.length > 0 && (
                <section className="page-section">
                    <h2>Deadlines</h2>
                    <ul className="result-list">
                        {results.deadlines.map(d => (
                            <li key={d.id}>
                                <button type="button" className="link-button" onClick={() => navigate("/deadlines")}>
                                    {d.title} · {relativeTime(d.date)}
                                </button>
                            </li>
                        ))}
                    </ul>
                </section>
            )}
        </PageLayout>
    );
};

export default SearchPage;
