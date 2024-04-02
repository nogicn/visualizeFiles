
## How to Use

1. **Prepare Input File**: Place a file in the `inputFiles` directory. This file can be of any type or content.

2. **Adjust Target Average**: In the `main.go` file, navigate to line 93. Adjust the `target average` until you achieve a result that isn't just noise. The value can range from 0.05 to 500 if necessary, though it typically falls between 0.1 and 404.

3. **Run the Program**:
   ```bash
   go run main.go
   ```

4. **Review Output**:
   - Output files will be generated in the `outputFiles` directory.
   - Refer to the example files provided in `outputFiles` to understand the output format and structure.
