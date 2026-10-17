import { afterEach, describe, expect, it, vi } from "vitest";
import { ApiError, apiRequest, setAuthToken, setUnauthorizedHandler } from "./client.ts";

const json = (body: unknown, status = 200) =>
    new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });

describe("apiRequest", () => {
    afterEach(() => {
        setAuthToken(null);
        setUnauthorizedHandler(null);
    });

    it("sends JSON with the token and returns the body", async () => {
        setAuthToken("abc");
        const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(json({ task: { id: "1" } }, 201));

        const result = await apiRequest<{ task: { id: string } }>("/tasks", { method: "POST", body: { title: "x" } });

        expect(result.task.id).toBe("1");
        const [url, init] = fetchMock.mock.calls[0];
        expect(url).toBe("http://localhost:8080/tasks");
        expect(init?.method).toBe("POST");
        expect(init?.body).toBe('{"title":"x"}');
        expect(init?.headers).toEqual({ "Content-Type": "application/json", Authorization: "Bearer abc" });
    });

    it("returns undefined for 204 No Content", async () => {
        vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 204 }));
        await expect(apiRequest("/tasks/1", { method: "DELETE" })).resolves.toBeUndefined();
    });

    it("throws ApiError with the server's message and field", async () => {
        vi.spyOn(globalThis, "fetch").mockResolvedValue(json({ error: "title is required", field: "title" }, 400));
        const err = await apiRequest("/tasks").catch(e => e);
        expect(err).toBeInstanceOf(ApiError);
        expect(err).toMatchObject({ status: 400, message: "title is required", field: "title" });
    });

    it("calls the unauthorized handler on 401 when logged in", async () => {
        setAuthToken("expired");
        const handler = vi.fn();
        setUnauthorizedHandler(handler);
        vi.spyOn(globalThis, "fetch").mockResolvedValue(json({ error: "invalid or expired token" }, 401));
        await expect(apiRequest("/tasks")).rejects.toMatchObject({ status: 401 });
        expect(handler).toHaveBeenCalledOnce();
    });

    it("reports network failures as status 0", async () => {
        vi.spyOn(globalThis, "fetch").mockRejectedValue(new TypeError("Failed to fetch"));
        await expect(apiRequest("/tasks")).rejects.toMatchObject({ status: 0 });
    });
});
