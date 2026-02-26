import { query, mutation } from "./_generated/server";
import { v } from "convex/values";

// Quick reference for Go integration (see main.go):
// - todoActions:list         <- GET    /todos
// - todoActions:add          <- POST   /add
// - todoActions:setCompleted <- PATCH  /update/:id
// - todoActions:setBody      <- PATCH  /update/body/:id
// - todoActions:remove       <- DELETE /delete/:id
//
// Note: "todoList" must match your Convex table name in schema.ts.
export const list = query({
  args: {},
  handler: async (ctx) => {
    // "todoList" is the name of the table that is set in your convex dashboard
    // if there is no convex dashboard, you can create one by running "npx convex dev" in your terminal
    // the name of the table should match the name of the table in your convex dashboard
    // so if the convex dashboard has a table named "todos", you should change "todoList" to "todos"
    return await ctx.db.query("todoList").order("desc").collect(); // "todoList" is the name of the table that is set in your convex dashboard
  },
});

// mutations are used to modify the data in the database, 
// and they can also return data after modifying it. 
export const add = mutation({
  // args are the parameters that are passed to the mutation when it is called from the client.
  // in this case, the add mutation takes a single argument called "body" which is a string.
  args: {
    body: v.string(),
  },
  // the handler is the function that is called when the mutation is executed.
  // it takes the context (ctx) and the arguments (args) as parameters.
  // the context (ctx) is an object that contains information about the current request, 
  // such as the database connection and the user making the request.
  handler: async (ctx, args) => {
    const id = await ctx.db.insert("todoList", {
      body: args.body,
      completed: false,
    });
    return await ctx.db.get(id);
  },
});

// Mark an existing todo as completed.
export const setCompleted = mutation({
  args: {
    id: v.id("todoList"),
  },
  handler: async (ctx, args) => {
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    await ctx.db.patch(args.id, { completed: true });
    return await ctx.db.get(args.id);
  },
});

// Update only the body text of an existing todo.
export const setBody = mutation({
  args: {
    id: v.id("todoList"),
    body: v.string(),
  },
  handler: async (ctx, args) => {
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    await ctx.db.patch(args.id, { body: args.body });
    return await ctx.db.get(args.id);
  },
});

// Delete a todo by id.
export const remove = mutation({
  args: {
    id: v.id("todoList"),
  },
  handler: async (ctx, args) => {
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    await ctx.db.delete(args.id);
  },
});

