const std = @import("std");
const builtin = @import("builtin");
const ArrayList = std.ArrayList;
const ChildProcess = std.ChildProcess;
const fmt = std.fmt;
const fs = std.fs;
const heap = std.heap;
const io = std.io;
const mem = std.mem;
const net = std.net;
const os = std.os;
const process = std.process;
const Thread = std.Thread;

const execName = "state-exec";

const Error = error{
    ArgIterator,
    ArgInvalidOne,
    ArgMissingOne,
    ArgInvalidTwo,
    ArgMissingTwo,
    ArgInvalidThree,
    ArgMissingThree,
    ArgInvalidFour,
    ArgMissingFour,
    ArgInvalidFive,
    ArgMissingFive,
    ArgCollector,
    ArgCollectRunt,
    ArgCollectUsr,
    InspectSelfPath,
    ThreadSpawn,
    ChildProcInit,
    ChildProcSpawn,
};

pub fn main() !void {
    const stderr = io.getStdErr().writer();

    const exitCode = run(stderr) catch |err| {
        try stderr.print("{s}: ", .{execName});

        switch (err) {
            Error.ArgIterator => try stderr.print("Cannot process args.\n", .{}),
            Error.ArgInvalidOne, Error.ArgMissingOne => try stderr.print("First arg should be a socket file.\n", .{}),
            Error.ArgInvalidTwo, Error.ArgMissingTwo => try stderr.print("Second arg should be a language runtime.\n", .{}),
            Error.ArgInvalidThree, Error.ArgMissingThree => try stderr.print("Third arg should be a project namespace.\n", .{}),
            Error.ArgInvalidFour, Error.ArgMissingFour => try stderr.print("Fourth arg should be a project commit ID.\n", .{}),
            Error.ArgInvalidFive, Error.ArgMissingFive => try stderr.print("Fifth arg should be a project headless status boolean.\n", .{}),
            Error.ArgCollector => try stderr.print("Cannot setup arg collector.\n", .{}),
            Error.ArgCollectRunt => try stderr.print("Cannot collect runtime arg.\n", .{}),
            Error.ArgCollectUsr => try stderr.print("Cannot collect user args.\n", .{}),
            Error.InspectSelfPath => try stderr.print("Cannot obtain path to this executable.\n", .{}),
            Error.ThreadSpawn => try stderr.print("Cannot spawn thread for heartbeat.\n", .{}),
            Error.ChildProcInit => try stderr.print("Cannot initialize child process for runtime.\n", .{}),
            Error.ChildProcSpawn => try stderr.print("Cannot spawn child process for runtime.\n", .{}),
        }

        try stderr.print("{s}: This application is not intended to be user serviceable; Please contact support for assistance.\n", .{execName});

        process.exit(1);
    };
    os.exit(exitCode);
}

const MsgData = struct {
    pub const fmt = "heart<{d}<{s}<{s}<{s}<{s}";

    pid: i32,
    exec: []const u8,
    nameSpace: []const u8,
    commitID: []const u8,
    headless: []const u8,
};

fn sendMsgToServer(a: mem.Allocator, stderr: fs.File.Writer, sock: []const u8, d: MsgData) !void {
    const conn = net.connectUnixSocket(sock) catch |err| {
        try stderr.print("{s}: Cannot connect to socket: {s}.\n", .{ execName, err });
        return;
    };
    defer conn.close();

    var clientMsg = try fmt.allocPrint(a, MsgData.fmt, .{ d.pid, d.exec, d.nameSpace, d.commitID, d.headless });
    _ = conn.write(clientMsg) catch |err| {
        try stderr.print("{s}: Cannot write to socket connection: {s}.\n", .{ execName, err });
        return;
    };

    var buf: [1024]u8 = undefined;
    _ = conn.read(buf[0..]) catch |err| {
        try stderr.print("{s}: Cannot read from socket connection: {s}.\n", .{ execName, err });
        return;
    };
}

fn run(stderr: fs.File.Writer) Error!u8 {
    var arena = heap.ArenaAllocator.init(heap.page_allocator);
    defer arena.deinit();
    const a = arena.allocator();

    var argIt = process.argsWithAllocator(a) catch return Error.ArgIterator;
    defer argIt.deinit();

    _ = argIt.skip();
    const sock = (argIt.next(a) orelse return Error.ArgMissingOne) catch return Error.ArgInvalidOne;
    const runt = (argIt.next(a) orelse return Error.ArgMissingTwo) catch return Error.ArgInvalidTwo;
    const nmsp = (argIt.next(a) orelse return Error.ArgMissingThree) catch return Error.ArgInvalidThree;
    const cmid = (argIt.next(a) orelse return Error.ArgMissingFour) catch return Error.ArgInvalidFour;
    const hdls = (argIt.next(a) orelse return Error.ArgMissingFive) catch return Error.ArgInvalidFive;

    var pid: i32 = @truncate(i32, @bitCast(i64, Thread.getCurrentId()));

    const exec = fs.selfExePathAlloc(a) catch return Error.InspectSelfPath;

    const data = MsgData{
        .pid = pid,
        .exec = exec,
        .nameSpace = nmsp,
        .commitID = cmid,
        .headless = hdls,
    };
    const clientThread = Thread.spawn(.{}, sendMsgToServer, .{ a, stderr, sock, data }) catch {
        return Error.ThreadSpawn;
    };
    defer clientThread.join();

    var usrArgs = process.argsAlloc(a) catch return Error.ArgCollector;
    defer process.argsFree(a, usrArgs);

    var cmdArgs = ArrayList([]const u8).init(a);
    defer cmdArgs.deinit();
    cmdArgs.append(runt) catch return Error.ArgCollectRunt;
    cmdArgs.appendSlice(usrArgs[6..]) catch return Error.ArgCollectUsr;

    const childProc = ChildProcess.init(cmdArgs.items, a) catch return Error.ChildProcInit;
    defer childProc.deinit();
    var term = childProc.spawnAndWait() catch return Error.ChildProcSpawn;
    return term.Exited;
}