#!/usr/bin/env python3
"""
Scheduler for the egg price data collector
Runs data collection at specified intervals
"""

import schedule
import time
import logging
from data_collector import EggPriceCollector

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

def run_collection():
    """Run the data collection process"""
    logger.info("Starting scheduled data collection")
    try:
        collector = EggPriceCollector()
        collector.run_collection()
        logger.info("Scheduled data collection completed")
    except Exception as e:
        logger.error(f"Scheduled data collection failed: {e}")

def main():
    """Main scheduler function"""
    logger.info("Starting egg price data collection scheduler")
    
    # Schedule data collection
    # Run daily at 8:00 AM
    schedule.every().day.at("08:00").do(run_collection)
    
    # Run weekly on Monday at 9:00 AM for comprehensive update
    schedule.every().monday.at("09:00").do(run_collection)
    
    # Run initial collection
    logger.info("Running initial data collection")
    run_collection()
    
    # Keep the scheduler running
    while True:
        schedule.run_pending()
        time.sleep(60)  # Check every minute

if __name__ == "__main__":
    main()
