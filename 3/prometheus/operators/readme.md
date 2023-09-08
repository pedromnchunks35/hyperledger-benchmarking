# Operators
- Most of the operators are the same in comparation to other programming languages so we will ignore the comparations and operations operators
## Logical/set binary operators
- and (intersection)
- or (union)
- unless (complement)

You can use those to aggregate vectors like so
```
vector1 and vector2
```
## Vector matching
- This is how we make one-to-many queries,many-to-one queries,left-hand queries and right-hand queries
### Vector matching keywords
- on
- ignoring
### Group modifiers
- group_left
- group_right
### One-to-on vector matches
- finds a unique pair of entries from each side of the operation
- ignoring keyword allows ignoring certain labels when matching
- on keyword allows reducing the set of considered labels to a provided list
#### Input example
```
method_code:http_errors:rate5m{method="get", code="500"}  24
method_code:http_errors:rate5m{method="get", code="404"}  30
method_code:http_errors:rate5m{method="put", code="501"}  3
method_code:http_errors:rate5m{method="post", code="500"} 6
method_code:http_errors:rate5m{method="post", code="404"} 21

method:http_requests:rate5m{method="get"}  600
method:http_requests:rate5m{method="del"}  34
method:http_requests:rate5m{method="post"} 120
```
#### Example query
```
method_code:http_errors:rate5m{code="500"} / ignoring(code) method:http_requests:rate5m
```
- this results in
```
{method="get"}  0.04            //  24 / 600
{method="post"} 0.05            //   6 / 120
```
- Without the ignoring, there would be no match.. we needed to remove the fields that should not match in order to make the this query using ignore
- So, by ignoring the code that isnt present in the requests,because only the method is, we could get the result
### Many-to-one and one-to-many vector matches
- In this perspective we use "group_left" or "group_right" modifiers
- Applying left or right determine which vector has the higher importance to become many for getting one-to-many or many-to-one
- one-to-many says that the right has more, group_right
- many-to-one says that the left has more, group_left
```
method_code:http_errors:rate5m / ignoring(code) group_left method:http_requests:rate5m
```
- The result of this query is
  ```
    {method="get", code="500"}  0.04            //  24 / 600
    {method="get", code="404"}  0.05            //  30 / 600
    {method="post", code="500"} 0.05            //   6 / 120
    {method="post", code="404"} 0.175           //  21 / 120
  ```
- By using group_left we are saying that the elements of the right side will be matched with multiple elements from the left.. so it is a many-to-one relationship
## Aggregation operators
|Property|Definition|
|--|--|
|sum|calculate sum over dimensions|
|min|calculate minimum over dimensions|
|max|calculate maximum over dimensions|
|avg|calculate the average over dimensions|
|group|all values in the resulting vector are 1|
|stddev|calculate population standard deviation over dimensions|
|stdvar|calculate population standard variance over dimensions|
|count|count number of elements in the vector|
|count_values|count number of elements with the same value|
|bottomk|smallest k elements by sample value|
|topk|largest k elements by sample value|
|quantile|calculate φ-quantile (0 ≤ φ ≤ 1) over dimensions (percentage over a dimension)|
- These aggregations can be used over all dimensions or by distinct dimensions by including **without** or **by** clause
### Syntax
```
<aggr-op> [without|by (<label list>)] ([parameter,] <vector expression>)
```
or
```
<aggr-op>([parameter,] <vector expression>) [without|by (<label list>)]
```
- **without** removes labels from the result set that are not on the without options
- **by** removes labels that are not in the by option
- **parameter** is a argument for the function of the operators that require aditional args like quantile,count_values,topk and bottomk
### Examples
```
 sum without (instance) (http_requests_total)
 sum by (application, group) (http_requests_total)
 sum(http_requests_total)
 count_values("version", build_version)
 topk(5, http_requests_total)
```