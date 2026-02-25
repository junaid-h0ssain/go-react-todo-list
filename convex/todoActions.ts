import { query, mutation } from "./_generated/server";
import { v } from "convex/values";


// queries are used to fetch data from the database, and they can also return data after fetching it.
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
    // in this insert call a new todo item is inserted into the "todoList" table 
    // with the body from the arguments and completed set to false.
    const id = await ctx.db.insert("todoList", {
      body: args.body,
      completed: false,
    });
    // after inserting the new todo item, the mutation returns the newly created todo item 
    // by fetching it from the database using its id.
    return await ctx.db.get(id);
  },
});

export const setCompleted = mutation({
  args: {
    // the id argument is defined as an id of the "todoList" table, 
    // which means that it should be a valid id of an item in the "todoList" table.
    id: v.id("todoList"),
  },
  handler: async (ctx, args) => {
    // the handler first checks if a todo item with the given id exists in the database.
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    // if the todo item exists, it updates the completed field of the todo item to true using the patch method.
    await ctx.db.patch(args.id, { completed: true });
    return await ctx.db.get(args.id);
  },
});

export const remove = mutation({
  args: {
    // same as before, this line defines the id argument as an id of the "todoList" table, 
    // which means that it should be a valid id of an item in the "todoList" table.
    id: v.id("todoList"),
  },
  handler: async (ctx, args) => {
    const todo = await ctx.db.get(args.id);
    if (!todo) {
      throw new Error("Todo not found");
    }
    // if the todo item exists, it deletes the todo item from the database using the delete method.
    await ctx.db.delete(args.id);
  },
});

