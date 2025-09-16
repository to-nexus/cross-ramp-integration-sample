using SampleGameBackend;

namespace ChecksumTestConsole
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("C# EVM Address Validation Test Console");
            Console.WriteLine("=====================================\n");

            try
            {
                // Run all tests
                ChecksumTest.RunTests();

                Console.WriteLine("\n🎉 All tests completed successfully!");
            }
            catch (Exception ex)
            {
                Console.WriteLine($"\n❌ Error running tests: {ex.Message}");
                Console.WriteLine($"Stack trace: {ex.StackTrace}");
                Environment.Exit(1);
            }

            Console.WriteLine("\nPress any key to exit...");
            Console.ReadKey();
        }
    }
}
