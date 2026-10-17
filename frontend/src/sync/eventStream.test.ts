import { afterEach, describe, expect, it, vi } from "vitest";
import { SSEParser, backoffDelay, connectEventStream, type StreamHandlers } from "./eventStream.ts";
import { setAuthToken } from "../api/client.ts";

describe("SSEParser", () => {
    it("parses complete messages and named events", () => {
        const parser = new SSEParser();
        const messages = parser.push('event: ready\ndata: {}\n\ndata: {"type":"task.deleted","id":"1"}\n\n');
        expect(messages).toEqual([
            { event: "ready", data: "{}" },
            { event: "message", data: '{"type":"task.deleted","id":"1"}' },
        ]);
    });

    it("handles chunks split in the middle of a line", () => {
        const parser = new SSEParser();
        expect(parser.push("da")).toEqual([]);
        expect(parser.push("ta: hel")).toEqual([]);
        expect(parser.push("lo\n")).toEqual([]);
        expect(parser.push("\n")).toEqual([{ event: "message", data: "hello" }]);
    });

    it("ignores keep-alive comments and joins multi-line data", () => {
        const parser = new SSEParser();
        const messages = parser.push(": ping\n\ndata: a\ndata: b\r\n\r\n");
        expect(messages).toEqual([{ event: "message", data: "a\nb" }]);
    });

    it("resets the event name after each message", () => {
        const parser = new SSEParser();
        const messages = parser.push("event: ready\ndata: 1\n\ndata: 2\n\n");
        expect(messages.map(m => m.event)).toEqual(["ready", "message"]);
    });
});

describe("backoffDelay", () => {
    it("doubles up to 30 seconds", () => {
        expect([0, 1, 2, 3, 10].map(backoffDelay)).toEqual([1000, 2000, 4000, 8000, 30000]);
    });
});

/** A fetch Response whose body streams the given chunks. */
const streamResponse = (chunks: string[], status = 200) =>
    new Response(
        new ReadableStream({
            start(controller) {
                for (const chunk of chunks) controller.enqueue(new TextEncoder().encode(chunk));
                controller.close();
            },
        }),
        { status },
    );

const handlers = (): StreamHandlers => ({
    onEvent: vi.fn(),
    onReady: vi.fn(),
    onStatus: vi.fn(),
    onUnauthorized: vi.fn(),
});

describe("connectEventStream", () => {
    afterEach(() => {
        setAuthToken(null);
        vi.useRealTimers();
    });

    it("sends the token and delivers ready and change events", async () => {
        setAuthToken("secret");
        const fetchMock = vi.spyOn(globalThis, "fetch").mockResolvedValue(
            streamResponse(["event: ready\ndata: {}\n\n", 'data: {"type":"task.deleted","id":"7"}\n\n']),
        );
        const h = handlers();
        const stop = connectEventStream(h);

        await vi.waitFor(() => expect(h.onEvent).toHaveBeenCalledWith({ type: "task.deleted", id: "7" }));
        expect(h.onReady).toHaveBeenCalledTimes(1);
        expect(h.onStatus).toHaveBeenCalledWith("live");
        const [url, init] = fetchMock.mock.calls[0];
        expect(String(url)).toMatch(/\/events$/);
        expect((init?.headers as Record<string, string>).Authorization).toBe("Bearer secret");
        stop();
    });

    it("reconnects after the stream ends", async () => {
        vi.useFakeTimers();
        const fetchMock = vi.spyOn(globalThis, "fetch").mockImplementation(async () => streamResponse(["event: ready\ndata: {}\n\n"]));
        const h = handlers();
        const stop = connectEventStream(h);

        await vi.waitFor(() => expect(h.onStatus).toHaveBeenCalledWith("offline"));
        await vi.advanceTimersByTimeAsync(1000);
        await vi.waitFor(() => expect(fetchMock).toHaveBeenCalledTimes(2));
        expect(h.onReady).toHaveBeenCalledTimes(2);
        stop();
    });

    it("stops and reports when the token is rejected", async () => {
        vi.spyOn(globalThis, "fetch").mockResolvedValue(new Response(null, { status: 401 }));
        const h = handlers();
        connectEventStream(h);
        await vi.waitFor(() => expect(h.onUnauthorized).toHaveBeenCalled());
        expect(h.onStatus).not.toHaveBeenCalledWith("offline");
    });
});
